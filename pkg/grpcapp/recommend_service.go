package grpcapp

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	randInt "math/rand"

	"github.com/thteam47/common-libs/ellipticutil"
	"github.com/thteam47/common-libs/x509util"
	recommendApi "github.com/thteam47/common/pkg/apiswagger/recommend-api"

	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/entity"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/pkg/entityutil"
	"github.com/thteam47/go-agent-pc/errutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *AgentpcService) RequestGenarateRecommendCbf(ctx context.Context, req *pb.StringRequest) (*pb.GenarateRecommendResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Ctx.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	customerID, err := entityutil.GetUserId(userContext)
	if err != nil {
		return nil, errutil.Wrap(err, "entityutil.GetUserId")
	}
	if customerID != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Ctx.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedData")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	recommendCbf := map[string]recommendApi.RecommendApiRecommend{}
	keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Value:    req.Ctx.DomainId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
			entity.FindRequestFilter{
				Key:      "UserId",
				Value:    req.Value,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "KeyInfoRepository.FindLast")
	}
	timeStartRequest := time.Now().UnixMilli()
	publicKeyX, err := x509util.ExtractKeyPublic(keyInfo.KeyPublic)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priX, err := x509util.ExtractKeyPrivate(keyInfo.KeyPrivate)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPrivate")
	}
	priXByte := new(big.Int).SetBytes(priX.D.Bytes())
	for j := 1; j <= int(combinedData.NumberItem1+combinedData.NumberItem2); j++ {
		cj, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, errutil.Wrap(err, "ecdsa.GenerateKey")
		}
		cjByte := new(big.Int).SetBytes(cj.D.Bytes())
		Cj2x, Cj2y := curve.ScalarBaseMult(cjByte.Bytes())
		rating := randInt.Intn(5)
		rij := new(big.Int).SetInt64(int64(rating))
		rijGx, rijGy := curve.ScalarBaseMult(rij.Bytes())

		cjXx, cjXy := curve.ScalarMult(publicKeyX.X, publicKeyX.Y, cjByte.Bytes())
		tongX, tongY := ellipticutil.AddPoint(curve, rijGx, rijGy, cjXx, cjXy)
		Cj1Decypt := elliptic.Marshal(curve, tongX, tongY)
		Cj2Decypt := elliptic.Marshal(curve, Cj2x, Cj2y)

		recommendCbf[strconv.Itoa(j)] = recommendApi.RecommendApiRecommend{
			ProcessDataC1: hex.EncodeToString(Cj1Decypt),
			ProcessDataC2: hex.EncodeToString(Cj2Decypt),
		}
	}

	timeEndRequest := time.Now().UnixMilli()
	dentaTimeRequest := timeEndRequest - timeStartRequest
	fmt.Printf("Time user CBF request: %d \n", dentaTimeRequest)

	timeStartServer := time.Now().UnixMilli()
	recommendCbfResult12, recommendCbfResult34, err := inst.componentsContainer.RecommendService().RecommendCbfGenarate(userContext, req.Ctx.DomainId, req.Ctx.TokenAgent, req.Value, recommendCbf)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.RecommendCbfGenarate")
	}

	timeEndServer := time.Now().UnixMilli()
	dentaTimeServer := timeEndServer - timeStartServer
	fmt.Printf("Time user CBF server: %d \n", dentaTimeServer)
	if recommendCbfResult12 == nil || len(recommendCbfResult12) != int(combinedData.NumberItem1+combinedData.NumberItem2) || recommendCbfResult34 == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "recommend Cbf Result bad request")
	}

	resultCk := map[int32]string{}
	timeStartLogarit := time.Now().UnixMilli()
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		processDataFk1 := ""
		processDataFk2 := ""
		if dataRecommendCbf, found := recommendCbfResult12[strconv.Itoa(k)]; found {
			processDataFk1 = dataRecommendCbf.ProcessDataFk1
			processDataFk2 = dataRecommendCbf.ProcessDataFk2
		}
		pointFk1, err := hex.DecodeString(processDataFk1)
		if err != nil {
			panic(err)
		}
		Fk1X, Fk1Y := elliptic.Unmarshal(curve, pointFk1)
		pointFk2, err := hex.DecodeString(processDataFk2)
		if err != nil {
			panic(err)
		}
		Fk2X, Fk2Y := elliptic.Unmarshal(curve, pointFk2)
		xF2kx, xF2ky := curve.ScalarMult(Fk2X, Fk2Y, priXByte.Bytes())

		// invert one point
		xF2kyY := new(big.Int).Neg(xF2ky)
		// point normalization
		xF2kyYSub := new(big.Int).Mod(xF2kyY, curve.Params().P)
		Ck3x, Ck3y := ellipticutil.AddPoint(curve, Fk1X, Fk1Y, xF2kx, xF2kyYSub)

		Ck3Decypt := elliptic.Marshal(curve, Ck3x, Ck3y)

		resultCk[int32(k)] = hex.EncodeToString(Ck3Decypt)
	}

	processDataFk3 := recommendCbfResult34.ProcessDataFk3
	processDataFk4 := recommendCbfResult34.ProcessDataFk4
	pointFk3, err := hex.DecodeString(processDataFk3)
	if err != nil {
		panic(err)
	}
	Fk3X, Fk3Y := elliptic.Unmarshal(curve, pointFk3)
	pointFk4, err := hex.DecodeString(processDataFk4)
	if err != nil {
		panic(err)
	}
	Fk4X, Fk4Y := elliptic.Unmarshal(curve, pointFk4)
	// C4
	xF4kx, xF4ky := curve.ScalarMult(Fk4X, Fk4Y, priXByte.Bytes())
	// invert one point
	xF4kyY := new(big.Int).Neg(xF4ky)
	// point normalization
	xF4kyYSub := new(big.Int).Mod(xF4kyY, curve.Params().P)
	Ck4x, Ck4y := ellipticutil.AddPoint(curve, Fk3X, Fk3Y, xF4kx, xF4kyYSub)

	maxCbfP := 5 * int(combinedData.NumberUser) * 5 * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100

	resultCkRo := map[int]int{}

	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		dk4 := 0

		for x := 1; x <= maxCbfP; x++ {
			xByte := new(big.Int).SetInt64(int64(x))
			xX, xY := curve.ScalarBaseMult(xByte.Bytes())
			point, err := hex.DecodeString(resultCk[int32(k)])
			if err != nil {
				panic(err)
			}
			AijX, AijY := elliptic.Unmarshal(curve, point)

			if xX != nil && xY != nil && AijX != nil && AijY != nil {
				if xX.Cmp(AijX) == 0 && xY.Cmp(AijY) == 0 {
					dk4 = x
				}
			}
		}
		resultCkRo[k] = dk4
	}

	d := 0
	maxCbfPD := 5 * int(combinedData.NumberUser) * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100
	for x := 1; x <= maxCbfPD; x++ {
		xByte := new(big.Int).SetInt64(int64(x))
		xX, xY := curve.ScalarBaseMult(xByte.Bytes())
		if xX != nil && xY != nil && Ck4x != nil && Ck4y != nil {
			if xX.Cmp(Ck4x) == 0 && xY.Cmp(Ck4y) == 0 {
				d = x
			}
		}
	}

	timeEndLogarit := time.Now().UnixMilli()
	dentaTimeLogarit := timeEndLogarit - timeStartLogarit
	fmt.Printf("Time user CBF Logarit: %d \n", dentaTimeLogarit)

	resultRespCbfData := []*pb.ResultRecommend{}
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		cbfItem := float64(resultCkRo[k]) / float64(d)
		numStr, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cbfItem), 64)
		if err != nil {
			return nil, errutil.Wrap(err, "strconv.ParseFloat")
		}
		resultRespCbfData = append(resultRespCbfData, &pb.ResultRecommend{
			PositionItem: int32(k),
			ProcessData:  float32(numStr),
		})
	}

	return &pb.GenarateRecommendResponse{
		Data: resultRespCbfData,
	}, nil
}

func (inst *AgentpcService) RequestGenarateRecommendCf(ctx context.Context, req *pb.StringRequest) (*pb.GenarateRecommendResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Ctx.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	customerID, err := entityutil.GetUserId(userContext)
	if err != nil {
		return nil, errutil.Wrap(err, "entityutil.GetUserId")
	}
	if customerID != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Ctx.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedData")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	recommendCf := map[string]recommendApi.RecommendApiRecommend{}
	keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Value:    req.Ctx.DomainId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
			entity.FindRequestFilter{
				Key:      "UserId",
				Value:    req.Value,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "KeyInfoRepository.FindLast")
	}
	timeStartRequest := time.Now().UnixMilli()
	publicKeyX, err := x509util.ExtractKeyPublic(keyInfo.KeyPublic)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priX, err := x509util.ExtractKeyPrivate(keyInfo.KeyPrivate)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priXByte := new(big.Int).SetBytes(priX.D.Bytes())
	for j := 1; j <= int(combinedData.NumberItem1+combinedData.NumberItem2); j++ {
		cj, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, errutil.Wrap(err, "ecdsa.GenerateKey")
		}
		cjByte := new(big.Int).SetBytes(cj.D.Bytes())
		Cj2x, Cj2y := curve.ScalarBaseMult(cjByte.Bytes())
		rating := randInt.Intn(5) + 1
		rij := new(big.Int).SetInt64(int64(rating * 10))
		rijGx, rijGy := curve.ScalarBaseMult(rij.Bytes())

		cjXx, cjXy := curve.ScalarMult(publicKeyX.X, publicKeyX.Y, cjByte.Bytes())
		tongX, tongY := ellipticutil.AddPoint(curve, rijGx, rijGy, cjXx, cjXy)
		Cj1Decypt := elliptic.Marshal(curve, tongX, tongY)
		Cj2Decypt := elliptic.Marshal(curve, Cj2x, Cj2y)

		recommendCf[strconv.Itoa(j)] = recommendApi.RecommendApiRecommend{
			ProcessDataC1: hex.EncodeToString(Cj1Decypt),
			ProcessDataC2: hex.EncodeToString(Cj2Decypt),
		}
	}
	timeEndRequest := time.Now().UnixMilli()
	dentaTimeRequest := timeEndRequest - timeStartRequest
	fmt.Printf("Time user CF request: %d \n", dentaTimeRequest)

	timeStartServer := time.Now().UnixMilli()
	recommendCfResult910, recommendCfResult1112, err := inst.componentsContainer.RecommendService().RecommendCfGenarate(userContext, req.Ctx.DomainId, req.Ctx.TokenAgent, req.Value, recommendCf)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.RecommendCbfGenarate")
	}

	timeEndServer := time.Now().UnixMilli()
	dentaTimeServer := timeEndServer - timeStartServer
	fmt.Printf("Time user CF server: %d \n", dentaTimeServer)

	if recommendCfResult910 == nil || len(recommendCfResult910) != int(combinedData.NumberItem1+combinedData.NumberItem2) || recommendCfResult1112 == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "recommend Cbf Result bad request")
	}

	resultCk := map[int32]string{}
	timeStartLogarit := time.Now().UnixMilli()
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		processDataFk9 := ""
		processDataFk10 := ""
		if dataRecommendCf, found := recommendCfResult910[strconv.Itoa(k)]; found {
			processDataFk9 = dataRecommendCf.ProcessDataFk9
			processDataFk10 = dataRecommendCf.ProcessDataFk10
		}
		pointFk9, err := hex.DecodeString(processDataFk9)
		if err != nil {
			return nil, errutil.Wrap(err, "hex.DecodeString")
		}
		Fk9X, Fk9Y := elliptic.Unmarshal(curve, pointFk9)
		pointFk10, err := hex.DecodeString(processDataFk10)
		if err != nil {
			return nil, errutil.Wrap(err, "hex.DecodeString")
		}
		Fk10X, Fk10Y := elliptic.Unmarshal(curve, pointFk10)

		xF10kx, xF10ky := curve.ScalarMult(Fk10X, Fk10Y, priXByte.Bytes())

		// invert one point
		xF10kyY := new(big.Int).Neg(xF10ky)
		// point normalization
		xF10kyYSub := new(big.Int).Mod(xF10kyY, curve.Params().P)
		Ck5x, Ck5y := ellipticutil.AddPoint(curve, Fk9X, Fk9Y, xF10kx, xF10kyYSub)

		Ck5Decypt := elliptic.Marshal(curve, Ck5x, Ck5y)

		resultCk[int32(k)] = hex.EncodeToString(Ck5Decypt)
	}

	processDataFk11 := recommendCfResult1112.ProcessDataFk11
	processDataFk12 := recommendCfResult1112.ProcessDataFk12
	pointF12, err := hex.DecodeString(processDataFk12)
	if err != nil {
		return nil, errutil.Wrap(err, "hex.DecodeString")
	}
	pointFk11, err := hex.DecodeString(processDataFk11)
	if err != nil {
		return nil, errutil.Wrap(err, "hex.DecodeString")
	}
	F12X, F12Y := elliptic.Unmarshal(curve, pointF12)
	xF12kx, xF12ky := curve.ScalarMult(F12X, F12Y, priXByte.Bytes())
	// invert one point
	xF12kyY := new(big.Int).Neg(xF12ky)
	// point normalization
	xF12kyYSub := new(big.Int).Mod(xF12kyY, curve.Params().P)
	Fk11X, Fk11Y := elliptic.Unmarshal(curve, pointFk11)
	//C6

	C6x, C6y := ellipticutil.AddPoint(curve, Fk11X, Fk11Y, xF12kx, xF12kyYSub)

	maxCbfP := 5 * int(combinedData.NumberUser) * 100 * int(combinedData.NumberItem1+combinedData.NumberItem2) * 10

	resultCkRo := map[int]int{}

	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		dk5 := 0

		for x := 1; x <= maxCbfP; x++ {
			xByte := new(big.Int).SetInt64(int64(x))
			xX, xY := curve.ScalarBaseMult(xByte.Bytes())
			point, err := hex.DecodeString(resultCk[int32(k)])
			if err != nil {
				return nil, errutil.Wrap(err, "hex.DecodeString")
			}
			AijX, AijY := elliptic.Unmarshal(curve, point)

			if xX != nil && xY != nil && AijX != nil && AijY != nil {
				if xX.Cmp(AijX) == 0 && xY.Cmp(AijY) == 0 {
					dk5 = x
					break
				}
			}
		}
		resultCkRo[k] = dk5
	}

	d := 0
	maxCbfPD := 5 * int(combinedData.NumberUser) * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100
	for x := 1; x <= maxCbfPD; x++ {
		xByte := new(big.Int).SetInt64(int64(x))
		xX, xY := curve.ScalarBaseMult(xByte.Bytes())
		if xX != nil && xY != nil && C6x != nil && C6y != nil {
			if xX.Cmp(C6x) == 0 && xY.Cmp(C6y) == 0 {
				d = x
				break
			}
		}
	}
	timeEndLogarit := time.Now().UnixMilli()
	dentaTimeLogarit := timeEndLogarit - timeStartLogarit
	fmt.Printf("Time user CF Logarit: %d \n", dentaTimeLogarit)

	resultRespCfData := []*pb.ResultRecommend{}
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		cbfItem := float64(resultCkRo[k]) / float64(d)
		numStr, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cbfItem), 64)
		if err != nil {
			return nil, errutil.Wrap(err, "strconv.ParseFloat")
		}
		resultRespCfData = append(resultRespCfData, &pb.ResultRecommend{
			PositionItem: int32(k),
			ProcessData:  float32(numStr),
		})
	}

	return &pb.GenarateRecommendResponse{
		Data: resultRespCfData,
	}, nil
}

func (inst *AgentpcService) RequestGenarateRecommendUserCbf(ctx context.Context, req *pb.RequestGenarateRecommendUser) (*pb.ListCategoryRecommendResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	if req.ProcessData == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Data Bad request")
	}
	userContext := entity.NewUserContext("default")
	userInfo := inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, req.Ctx.TokenAgent)
	if userInfo == nil {
		return nil, status.Errorf(codes.Unauthenticated, "userInfo Unauthenticated")
	}
	accessToken := inst.componentsContainer.IdentityAuthenService().AccessToken(userContext, req.Ctx.TokenAgent)

	if accessToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "accessToken Unauthenticated")
	}
	userContext.SetAccessToken(accessToken)
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedData")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	recommendCbf := map[string]recommendApi.RecommendApiRecommend{}
	keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Value:    userInfo.DomainId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
			entity.FindRequestFilter{
				Key:      "UserId",
				Value:    userInfo.UserId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "KeyInfoRepository.FindLast")
	}
	timeStartRequest := time.Now().UnixMilli()
	publicKeyX, err := x509util.ExtractKeyPublic(keyInfo.KeyPublic)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priX, err := x509util.ExtractKeyPrivate(keyInfo.KeyPrivate)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPrivate")
	}
	priXByte := new(big.Int).SetBytes(priX.D.Bytes())
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for j := 1; j <= int(combinedData.NumberItem1+combinedData.NumberItem2); j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			cj, err := ecdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				log.Println(errutil.Wrap(err, "ecdsa.GenerateKey"))
			}
			cjByte := new(big.Int).SetBytes(cj.D.Bytes())
			Cj2x, Cj2y := curve.ScalarBaseMult(cjByte.Bytes())
			rating := 0
			if ratingTmp, found := req.ProcessData[strconv.Itoa(j)]; found {
				rating = int(ratingTmp)
			}
			rij := new(big.Int).SetInt64(int64(rating))
			rijGx, rijGy := curve.ScalarBaseMult(rij.Bytes())

			cjXx, cjXy := curve.ScalarMult(publicKeyX.X, publicKeyX.Y, cjByte.Bytes())
			tongX, tongY := ellipticutil.AddPoint(curve, rijGx, rijGy, cjXx, cjXy)
			Cj1Decypt := elliptic.Marshal(curve, tongX, tongY)
			Cj2Decypt := elliptic.Marshal(curve, Cj2x, Cj2y)

			mutex.Lock()
			recommendCbf[strconv.Itoa(j)] = recommendApi.RecommendApiRecommend{
				ProcessDataC1: hex.EncodeToString(Cj1Decypt),
				ProcessDataC2: hex.EncodeToString(Cj2Decypt),
			}
			mutex.Unlock()
		}(j)
	}

	wg.Wait()

	timeEndRequest := time.Now().UnixMilli()
	dentaTimeRequest := timeEndRequest - timeStartRequest
	fmt.Printf("Time user CBF request: %d \n", dentaTimeRequest)

	timeStartServer := time.Now().UnixMilli()
	recommendCbfResult12, recommendCbfResult34, err := inst.componentsContainer.RecommendService().RecommendCbfGenarate(userContext, userInfo.DomainId, req.Ctx.TokenAgent, userInfo.UserId, recommendCbf)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.RecommendCbfGenarate")
	}

	if recommendCbfResult12 == nil || len(recommendCbfResult12) != int(combinedData.NumberItem1+combinedData.NumberItem2) || recommendCbfResult34 == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "recommend Cbf Result bad request")
	}
	timeEndServer := time.Now().UnixMilli()
	dentaTimeServer := timeEndServer - timeStartServer
	fmt.Printf("Time user CBF server: %d \n", dentaTimeServer)
	resultCk := map[int32]string{}
	timeStartLogarit := time.Now().UnixMilli()
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			processDataFk1 := ""
			processDataFk2 := ""
			if dataRecommendCbf, found := recommendCbfResult12[strconv.Itoa(k)]; found {
				processDataFk1 = dataRecommendCbf.ProcessDataFk1
				processDataFk2 = dataRecommendCbf.ProcessDataFk2
			}
			pointFk1, err := hex.DecodeString(processDataFk1)
			if err != nil {
				panic(err)
			}
			Fk1X, Fk1Y := elliptic.Unmarshal(curve, pointFk1)
			pointFk2, err := hex.DecodeString(processDataFk2)
			if err != nil {
				panic(err)
			}
			Fk2X, Fk2Y := elliptic.Unmarshal(curve, pointFk2)
			xF2kx, xF2ky := curve.ScalarMult(Fk2X, Fk2Y, priXByte.Bytes())

			// invert one point
			xF2kyY := new(big.Int).Neg(xF2ky)
			// point normalization
			xF2kyYSub := new(big.Int).Mod(xF2kyY, curve.Params().P)
			Ck3x, Ck3y := ellipticutil.AddPoint(curve, Fk1X, Fk1Y, xF2kx, xF2kyYSub)

			Ck3Decypt := elliptic.Marshal(curve, Ck3x, Ck3y)

			mutex.Lock()
			resultCk[int32(k)] = hex.EncodeToString(Ck3Decypt)
			mutex.Unlock()
		}(k)
	}

	wg.Wait()

	processDataFk3 := recommendCbfResult34.ProcessDataFk3
	processDataFk4 := recommendCbfResult34.ProcessDataFk4
	pointFk3, err := hex.DecodeString(processDataFk3)
	if err != nil {
		panic(err)
	}
	Fk3X, Fk3Y := elliptic.Unmarshal(curve, pointFk3)
	pointFk4, err := hex.DecodeString(processDataFk4)
	if err != nil {
		panic(err)
	}
	Fk4X, Fk4Y := elliptic.Unmarshal(curve, pointFk4)
	// C4
	xF4kx, xF4ky := curve.ScalarMult(Fk4X, Fk4Y, priXByte.Bytes())
	// invert one point
	xF4kyY := new(big.Int).Neg(xF4ky)
	// point normalization
	xF4kyYSub := new(big.Int).Mod(xF4kyY, curve.Params().P)
	Ck4x, Ck4y := ellipticutil.AddPoint(curve, Fk3X, Fk3Y, xF4kx, xF4kyYSub)

	maxCbfP := 5 * int(combinedData.NumberUser) * 5 * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100

	resultCkRo := map[int]int{}

	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			dk4 := 0

			for x := 1; x <= maxCbfP; x++ {
				xByte := new(big.Int).SetInt64(int64(x))
				xX, xY := curve.ScalarBaseMult(xByte.Bytes())
				point, err := hex.DecodeString(resultCk[int32(k)])
				if err != nil {
					panic(err)
				}
				AijX, AijY := elliptic.Unmarshal(curve, point)

				if xX != nil && xY != nil && AijX != nil && AijY != nil {
					if xX.Cmp(AijX) == 0 && xY.Cmp(AijY) == 0 {
						dk4 = x
					}
				}
			}
			mutex.Lock()
			resultCkRo[k] = dk4
			mutex.Unlock()
		}(k)
	}

	wg.Wait()
	d := 0
	maxCbfPD := 5 * int(combinedData.NumberUser) * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100
	for x := 1; x <= maxCbfPD; x++ {
		xByte := new(big.Int).SetInt64(int64(x))
		xX, xY := curve.ScalarBaseMult(xByte.Bytes())
		if xX != nil && xY != nil && Ck4x != nil && Ck4y != nil {
			if xX.Cmp(Ck4x) == 0 && xY.Cmp(Ck4y) == 0 {
				d = x
			}
		}
	}
	timeEndLogarit := time.Now().UnixMilli()
	dentaTimeLogarit := timeEndLogarit - timeStartLogarit
	fmt.Printf("Time user CBF Logarit: %d \n", dentaTimeLogarit)

	resultRespCbfData := map[int32]float32{}
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		cbfItem := float64(resultCkRo[k]) / float64(d)
		numStr, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cbfItem), 64)
		if err != nil {
			return nil, errutil.Wrap(err, "strconv.ParseFloat")
		}
		resultRespCbfData[int32(k)] = float32(numStr)
	}

	categories, total, err := inst.componentsContainer.SurveyService().GetCategoriesByRecommendTenant(userContext.SetDomainId(userInfo.DomainId).Clone().EscalatePrivilege(), userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}

	items, err := makeCategoriesRecommend(categories, resultRespCbfData)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategoriesRecommend")
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProcessData > items[j].ProcessData
	})

	return &pb.ListCategoryRecommendResponse{
		Data:  items,
		Total: total,
	}, nil
}

func (inst *AgentpcService) RequestGenarateRecommendUserCf(ctx context.Context, req *pb.RequestGenarateRecommendUser) (*pb.ListCategoryRecommendResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	if req.ProcessData == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Data Bad request")
	}
	userContext := entity.NewUserContext("default")
	userInfo := inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, req.Ctx.TokenAgent)
	if userInfo == nil {
		return nil, status.Errorf(codes.Unauthenticated, "userInfo Unauthenticated")
	}
	accessToken := inst.componentsContainer.IdentityAuthenService().AccessToken(userContext, req.Ctx.TokenAgent)

	if accessToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "accessToken Unauthenticated")
	}
	userContext.SetAccessToken(accessToken)
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedData")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	recommendCf := map[string]recommendApi.RecommendApiRecommend{}
	keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Value:    userInfo.DomainId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
			entity.FindRequestFilter{
				Key:      "UserId",
				Value:    userInfo.UserId,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "KeyInfoRepository.FindLast")
	}
	timeStartRequest := time.Now().UnixMilli()
	publicKeyX, err := x509util.ExtractKeyPublic(keyInfo.KeyPublic)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priX, err := x509util.ExtractKeyPrivate(keyInfo.KeyPrivate)
	if err != nil {
		return nil, errutil.Wrap(err, "x509util.ExtractKeyPublic")
	}
	priXByte := new(big.Int).SetBytes(priX.D.Bytes())
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for j := 1; j <= int(combinedData.NumberItem1+combinedData.NumberItem2); j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			cj, err := ecdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				log.Println(errutil.Wrap(err, "ecdsa.GenerateKey").Error())
			}
			cjByte := new(big.Int).SetBytes(cj.D.Bytes())
			Cj2x, Cj2y := curve.ScalarBaseMult(cjByte.Bytes())
			rating := 0
			if ratingTmp, found := req.ProcessData[strconv.Itoa(j)]; found {
				rating = int(ratingTmp)
			}
			rij := new(big.Int).SetInt64(int64(rating * 10))
			rijGx, rijGy := curve.ScalarBaseMult(rij.Bytes())

			cjXx, cjXy := curve.ScalarMult(publicKeyX.X, publicKeyX.Y, cjByte.Bytes())
			tongX, tongY := ellipticutil.AddPoint(curve, rijGx, rijGy, cjXx, cjXy)
			Cj1Decypt := elliptic.Marshal(curve, tongX, tongY)
			Cj2Decypt := elliptic.Marshal(curve, Cj2x, Cj2y)
			mutex.Lock()
			recommendCf[strconv.Itoa(j)] = recommendApi.RecommendApiRecommend{
				ProcessDataC1: hex.EncodeToString(Cj1Decypt),
				ProcessDataC2: hex.EncodeToString(Cj2Decypt),
			}
			mutex.Unlock()
		}(j)
	}

	wg.Wait()
	timeEndRequest := time.Now().UnixMilli()
	dentaTimeRequest := timeEndRequest - timeStartRequest
	fmt.Printf("Time user CF request: %d \n", dentaTimeRequest)

	timeStartServer := time.Now().UnixMilli()
	recommendCfResult910, recommendCfResult1112, err := inst.componentsContainer.RecommendService().RecommendCfGenarate(userContext, userInfo.DomainId, req.Ctx.TokenAgent, userInfo.UserId, recommendCf)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.RecommendCbfGenarate")
	}

	if recommendCfResult910 == nil || len(recommendCfResult910) != int(combinedData.NumberItem1+combinedData.NumberItem2) || recommendCfResult1112 == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "recommend Cbf Result bad request")
	}
	timeEndServer := time.Now().UnixMilli()
	dentaTimeServer := timeEndServer - timeStartServer
	fmt.Printf("Time user CF server: %d \n", dentaTimeServer)
	resultCk := map[int32]string{}
	timeStartLogarit := time.Now().UnixMilli()
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		wg.Add(1)
		go func(k int) {
			wg.Done()
			processDataFk9 := ""
			processDataFk10 := ""
			if dataRecommendCf, found := recommendCfResult910[strconv.Itoa(k)]; found {
				processDataFk9 = dataRecommendCf.ProcessDataFk9
				processDataFk10 = dataRecommendCf.ProcessDataFk10
			}
			pointFk9, err := hex.DecodeString(processDataFk9)
			if err != nil {
				log.Println(errutil.Wrap(err, "hex.DecodeString").Error())
			}
			Fk9X, Fk9Y := elliptic.Unmarshal(curve, pointFk9)
			pointFk10, err := hex.DecodeString(processDataFk10)
			if err != nil {
				log.Println(errutil.Wrap(err, "hex.DecodeString").Error())
			}
			Fk10X, Fk10Y := elliptic.Unmarshal(curve, pointFk10)

			xF10kx, xF10ky := curve.ScalarMult(Fk10X, Fk10Y, priXByte.Bytes())

			// invert one point
			xF10kyY := new(big.Int).Neg(xF10ky)
			// point normalization
			xF10kyYSub := new(big.Int).Mod(xF10kyY, curve.Params().P)
			Ck5x, Ck5y := ellipticutil.AddPoint(curve, Fk9X, Fk9Y, xF10kx, xF10kyYSub)

			Ck5Decypt := elliptic.Marshal(curve, Ck5x, Ck5y)

			mutex.Lock()
			resultCk[int32(k)] = hex.EncodeToString(Ck5Decypt)
			mutex.Unlock()
		}(k)
	}

	wg.Wait()
	processDataFk11 := recommendCfResult1112.ProcessDataFk11
	processDataFk12 := recommendCfResult1112.ProcessDataFk12
	pointF12, err := hex.DecodeString(processDataFk12)
	if err != nil {
		return nil, errutil.Wrap(err, "hex.DecodeString")
	}
	pointFk11, err := hex.DecodeString(processDataFk11)
	if err != nil {
		return nil, errutil.Wrap(err, "hex.DecodeString")
	}
	F12X, F12Y := elliptic.Unmarshal(curve, pointF12)
	xF12kx, xF12ky := curve.ScalarMult(F12X, F12Y, priXByte.Bytes())
	// invert one point
	xF12kyY := new(big.Int).Neg(xF12ky)
	// point normalization
	xF12kyYSub := new(big.Int).Mod(xF12kyY, curve.Params().P)
	Fk11X, Fk11Y := elliptic.Unmarshal(curve, pointFk11)
	//C6

	C6x, C6y := ellipticutil.AddPoint(curve, Fk11X, Fk11Y, xF12kx, xF12kyYSub)

	maxCbfP := 5 * int(combinedData.NumberUser) * 100 * int(combinedData.NumberItem1+combinedData.NumberItem2) * 10

	resultCkRo := map[int]int{}

	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			dk5 := 0

			for x := 1; x <= maxCbfP; x++ {
				xByte := new(big.Int).SetInt64(int64(x))
				xX, xY := curve.ScalarBaseMult(xByte.Bytes())
				point, err := hex.DecodeString(resultCk[int32(k)])
				if err != nil {
					log.Println(errutil.Wrap(err, "hex.DecodeString"))
				}
				AijX, AijY := elliptic.Unmarshal(curve, point)

				if xX != nil && xY != nil && AijX != nil && AijY != nil {
					if xX.Cmp(AijX) == 0 && xY.Cmp(AijY) == 0 {
						dk5 = x
						break
					}
				}
			}
			mutex.Lock()
			resultCkRo[k] = dk5
			mutex.Unlock()
		}(k)
	}

	wg.Wait()
	d := 0
	maxCbfPD := 5 * int(combinedData.NumberUser) * int(combinedData.NumberItem1+combinedData.NumberItem2) * 100
	for x := 1; x <= maxCbfPD; x++ {
		xByte := new(big.Int).SetInt64(int64(x))
		xX, xY := curve.ScalarBaseMult(xByte.Bytes())
		if xX != nil && xY != nil && C6x != nil && C6y != nil {
			if xX.Cmp(C6x) == 0 && xY.Cmp(C6y) == 0 {
				d = x
				break
			}
		}
	}
	timeEndLogarit := time.Now().UnixMilli()
	dentaTimeLogarit := timeEndLogarit - timeStartLogarit
	fmt.Printf("Time user CF Logarit: %d \n", dentaTimeLogarit)
	resultRespCfData := map[int32]float32{}
	for k := 1; k <= int(combinedData.NumberItem1+combinedData.NumberItem2); k++ {
		cbfItem := float64(resultCkRo[k]) / float64(d)
		numStr, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cbfItem), 64)
		if err != nil {
			return nil, errutil.Wrap(err, "strconv.ParseFloat")
		}
		resultRespCfData[int32(k)] = float32(numStr)
	}

	categories, total, err := inst.componentsContainer.SurveyService().GetCategoriesByRecommendTenant(userContext.SetDomainId(userInfo.DomainId).Clone().EscalatePrivilege(), userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}

	items, err := makeCategoriesRecommend(categories, resultRespCfData)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategoriesRecommend")
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProcessData > items[j].ProcessData
	})

	return &pb.ListCategoryRecommendResponse{
		Data:  items,
		Total: total,
	}, nil
}
