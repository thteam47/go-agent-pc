package grpcapp

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	randInt "math/rand"
	"sync"
	"time"

	"github.com/thteam47/common-libs/ellipticutil"
	"github.com/thteam47/common-libs/x509util"
	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/entity"
	recommendApi "github.com/thteam47/common/pkg/apiswagger/recommend-api"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var curve = elliptic.P256()

func (inst *AgentpcService) FakeGenerateKeyInfo(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	totalTimeKeyUser := int64(0)
	fmt.Println("time.Now().UnixMilli()", time.Now().UnixMilli())
	coutKeyUser := 0
	totalTimeKeyItemOnePart := int64(0)
	coutKeyItemOnePart := 0
	totalTimeKeyItemTwoPart := int64(0)
	coutKeyItemTwoPart := 0
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
				Filters: []entity.FindRequestFilter{
					entity.FindRequestFilter{
						Key:      "DomainId",
						Value:    user.DomainId,
						Operator: entity.FindRequestFilterOperatorEqualTo,
					},
					entity.FindRequestFilter{
						Key:      "UserId",
						Value:    user.UserId,
						Operator: entity.FindRequestFilterOperatorEqualTo,
					},
				},
			})
			if err != nil {
				log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
				return
			}
			if keyInfo == nil {
				timeOneUserStart := time.Now().UnixMilli()
				privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
				if err != nil {
					panic(err)
				}
				publicKey := &privateKey.PublicKey
				privateKeyGen, err := x509util.GenerateKeyPrivate(privateKey)
				if err != nil {
					panic(err)
				}
				publicKeyGen, err := x509util.GenerateKeyPublic(publicKey)
				if err != nil {
					panic(err)
				}
				keyInfo = &models.KeyInfo{
					UserId:         user.UserId,
					DomainId:       user.DomainId,
					PositionUserId: user.Position,
					KeyPrivate:     privateKeyGen,
					KeyPublic:      publicKeyGen,
				}
				err = inst.componentsContainer.KeyInfoRepository().Create(userContext, keyInfo)
				if err != nil {
					log.Printf(errutil.Wrap(err, "KeyInfoRepository.Create").Error())
				}
				timeOneUserEnd := time.Now().UnixMilli()
				dentaTime := timeOneUserEnd - timeOneUserStart
				mutex.Lock()
				totalTimeKeyUser += dentaTime
				coutKeyUser++
				mutex.Unlock()
			}
			err = inst.componentsContainer.RecommendService().KeyPublicUserSend(userContext, req.Value, req.Ctx.TokenAgent, keyInfo)
			if err != nil {
				log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUserSend").Error())
			}
			timeOneUserTwoPartStart := time.Now().UnixMilli()
			for j := 1; j <= int(combinedData.NkTwoPart); j++ {
				keyInfoItem, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Value:    user.DomainId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Value:    user.UserId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Value:    j,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "Part",
							Value:    2,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
					},
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
				}
				if keyInfoItem == nil {
					privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
					if err != nil {
						panic(err)
					}
					publicKey := &privateKey.PublicKey
					privateKeyGen, err := x509util.GenerateKeyPrivate(privateKey)
					if err != nil {
						panic(err)
					}
					publicKeyGen, err := x509util.GenerateKeyPublic(publicKey)
					if err != nil {
						panic(err)
					}
					keyInfoItem = &models.KeyInfoItem{
						UserId:       user.UserId,
						DomainId:     user.DomainId,
						PositionItem: int32(j),
						PositionUser: user.Position,
						KeyPrivate:   privateKeyGen,
						KeyPublic:    publicKeyGen,
						Part:         2,
					}
					err = inst.componentsContainer.KeyInfoItemRepository().Create(userContext, keyInfoItem)
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.Create").Error())
					}
				}
				err = inst.componentsContainer.RecommendService().KeyPublicItemSend(userContext, req.Value, req.Ctx.TokenAgent, keyInfoItem)
				if err != nil {
					log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicItemSend").Error())
				}
			}
			timeOneUserTwoPartEnd := time.Now().UnixMilli()
			dentaTimeTwoPart := timeOneUserTwoPartEnd - timeOneUserTwoPartStart
			mutex.Lock()
			totalTimeKeyItemTwoPart += dentaTimeTwoPart
			coutKeyItemTwoPart++
			mutex.Unlock()
			nkOnePart := int(combinedData.NkOnePart1)
			if user.DomainId == combinedData.TenantId2 {
				nkOnePart = int(combinedData.NkOnePart2)
			}
			timeOneUserOnePartStart := time.Now().UnixMilli()
			for j := 1; j <= nkOnePart; j++ {
				keyInfoItem, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Value:    user.DomainId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Value:    user.UserId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Value:    j,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "Part",
							Value:    1,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
					},
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
				}
				if keyInfoItem == nil {
					privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
					if err != nil {
						panic(err)
					}
					publicKey := &privateKey.PublicKey
					privateKeyGen, err := x509util.GenerateKeyPrivate(privateKey)
					if err != nil {
						panic(err)
					}
					publicKeyGen, err := x509util.GenerateKeyPublic(publicKey)
					if err != nil {
						panic(err)
					}
					keyInfoItem = &models.KeyInfoItem{
						UserId:       user.UserId,
						DomainId:     user.DomainId,
						PositionItem: int32(j),
						PositionUser: user.Position,
						KeyPrivate:   privateKeyGen,
						KeyPublic:    publicKeyGen,
						Part:         1,
					}
					err = inst.componentsContainer.KeyInfoItemRepository().Create(userContext, keyInfoItem)
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.Create").Error())
					}
				}
				err = inst.componentsContainer.RecommendService().KeyPublicItemSend(userContext, req.Value, req.Ctx.TokenAgent, keyInfoItem)
				if err != nil {
					log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicItemSend").Error())
				}
			}
			timeOneUserOnePartEnd := time.Now().UnixMilli()
			dentaTimeOnePart := timeOneUserOnePartEnd - timeOneUserOnePartStart
			mutex.Lock()
			totalTimeKeyItemOnePart += dentaTimeOnePart
			coutKeyItemOnePart++
			mutex.Unlock()
		}(user)
	}

	wg.Wait()

	fmt.Printf("Time gen key %d user: %d \n", coutKeyUser, totalTimeKeyUser)
	fmt.Printf("Time gen key 1 user: %d \n", int(totalTimeKeyUser)/coutKeyUser)
	fmt.Printf("Time gen key item one part %d user: %d \n", coutKeyItemOnePart, totalTimeKeyItemOnePart)
	fmt.Printf("Time gen key item one part 1 user: %d \n", int(totalTimeKeyItemOnePart)/coutKeyItemOnePart)
	fmt.Printf("Time gen key item two part %d user: %d \n", coutKeyItemTwoPart, totalTimeKeyItemTwoPart)
	fmt.Printf("Time gen key item two part 1 user: %d \n", int(totalTimeKeyItemTwoPart)/coutKeyItemTwoPart)
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) FakeDoSurveyGenerate(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			surveys, err := inst.componentsContainer.SurveyService().GetSurveysByUserId(userContext.Clone().EscalatePrivilege(), user.DomainId, user.UserId)
			if err != nil {
				log.Printf(errutil.Wrap(err, "SurveyService.GetSurveysByUserId").Error())
			}
			for _, survey := range surveys {
				for _, question := range survey.Questions {
					rating := randInt.Intn(len(question.Answers))
					err = inst.componentsContainer.ResultCardRepository().Create(userContext, &models.ResultCard{
						DomainId:       user.DomainId,
						UserId:         user.UserId,
						PositionUser:   user.Position,
						CategoryId:     question.CategoryId,
						PositionItem:   question.Position,
						SurveyId:       survey.SurveyId,
						PositionOption: int32(rating) + 1,
						Option:         question.Answers[rating],
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "ResultCardRepository.Create").Error())
					}
				}
			}
		}(user)
	}
	wg.Wait()
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) FakeSendProcessDataSurveyOnePart(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	nkOnePart := combinedData.NkOnePart1
	if req.Value == combinedData.TenantId2 {
		nkOnePart = combinedData.NkOnePart2
	}
	sOnePart := combinedData.SOnePart1
	if req.Value == combinedData.TenantId2 {
		sOnePart = combinedData.SOnePart2
	}
	totalTimeSendData := int64(0)
	countUser := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			j := 1
			isStopj := false
			timeOneUserStart := time.Now().UnixMilli()
			for t := 1; t <= int(nkOnePart)-1; t++ {
				if isStopj {
					break
				}
				for k := t + 1; k <= int(nkOnePart); k++ {
					processDataSurvey, err := inst.componentsContainer.ProcessDataSurveyRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    user.DomainId,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    user.UserId,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    j,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "ProcessDataSurveyRepository.FindLast").Error())
					}
					processedData := int32(0)
					if processDataSurvey != nil {
						processedData = processDataSurvey.ProcessedData
					}
					aijByte := new(big.Int).SetInt64(int64(processedData))
					ajGx, ajGy := curve.ScalarBaseMult(aijByte.Bytes())

					//keyInfoIK
					keyInfoIK, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    k,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    1,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoItemRepository.FindLast").Error())
					}

					// keyInfoIT
					keyInfoIT, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    t,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    1,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}

					//keyPublicT
					keyPublicT, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(t), 1)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}
					//keyPublicK
					keyPublicK, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(k), 1)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}
					privateKeyIK, err := x509util.ExtractKeyPrivate(keyInfoIK.KeyPrivate)
					if err != nil {
						log.Printf(errutil.Wrap(err, "x509util.ExtractKeyPrivate").Error())
					}
					privateKeyIT, err := x509util.ExtractKeyPrivate(keyInfoIT.KeyPrivate)
					if err != nil {
						log.Printf(errutil.Wrap(err, "x509util.ExtractKeyPrivate").Error())
					}
					publicKeyT, err := x509util.ExtractKeyPublic(keyPublicT.KeyPublic)
					if err != nil {
						log.Printf(errutil.Wrap(err, "x509util.ExtractKeyPrivate").Error())
					}
					publicKeyK, err := x509util.ExtractKeyPublic(keyPublicK.KeyPublic)
					if err != nil {
						log.Printf(errutil.Wrap(err, "x509util.ExtractKeyPrivate").Error())
					}
					valueprivateKeyIK := new(big.Int).SetBytes(privateKeyIK.D.Bytes())
					kuikKptX, kuikKptY := curve.ScalarMult(publicKeyT.X, publicKeyT.Y, valueprivateKeyIK.Bytes())
					valueprivateKeyIT := new(big.Int).SetBytes(privateKeyIT.D.Bytes())
					kuitKpkX, kuikKpkY := curve.ScalarMult(publicKeyK.X, publicKeyK.Y, valueprivateKeyIT.Bytes())
					// invert one point
					kuikKpkYY := new(big.Int).Neg(kuikKpkY)
					// point normalization
					kuikKpkYSub := new(big.Int).Mod(kuikKpkYY, curve.Params().P)
					x3, y3 := ellipticutil.AddPoint(curve, kuikKptX, kuikKptY, kuitKpkX, kuikKpkYSub)
					tongX, tongY := ellipticutil.AddPoint(curve, ajGx, ajGy, x3, y3)
					decryptIJ := elliptic.Marshal(curve, tongX, tongY)

					err = inst.componentsContainer.RecommendService().ProcessedDataSendPart(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiResultCard{
						UserId:                user.UserId,
						PositionUser:          user.Position,
						PositionItem:          int32(j),
						ProcessedData:         hex.EncodeToString(decryptIJ),
						PositionItemOriginal:  processDataSurvey.PositionItemOriginal,
						PositionItemOriginal1: processDataSurvey.PositionItemOriginal1,
						PositionItemOriginal2: processDataSurvey.PositionItemOriginal2,
						Part:                  1,
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendPart").Error())
					}
					if j == int(sOnePart) {
						isStopj = true
						break
					} else {
						j++
					}
				}
			}
			timeOneUserEnd := time.Now().UnixMilli()
			dentaTime := timeOneUserEnd - timeOneUserStart
			mutex.Lock()
			totalTimeSendData += dentaTime
			countUser++
			mutex.Unlock()
		}(user)
	}
	wg.Wait()
	fmt.Printf("Time send data of %d user: %d \n", countUser, totalTimeSendData)
	fmt.Printf("Time send data 1 user: %d \n", int(totalTimeSendData)/countUser)
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) FakeSendProcessDataSurveyPhase3TwoPart(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	if combinedData.TenantId1 != req.Value {
		return nil, status.Errorf(codes.PermissionDenied, "Tenant is not allow create phase3 two part")
	}
	totalTimeSendDataPhase3 := int64(0)
	countUser := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			timeOneUserStart := time.Now().UnixMilli()
			for j := 1; j <= int(combinedData.STwoPart); j++ {
				keyInfoItemPhase3, err := inst.componentsContainer.KeyInfoItemPhase3Repository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Value:    user.DomainId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Value:    user.UserId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Value:    int32(j),
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
					},
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
				}
				if keyInfoItemPhase3 == nil {
					privateKeyCij, err := ecdsa.GenerateKey(curve, rand.Reader)
					if err != nil {
						panic(err)
					}
					publicKeyCij := &privateKeyCij.PublicKey
					privateKeyGen, err := x509util.GenerateKeyPrivate(privateKeyCij)
					if err != nil {
						panic(err)
					}
					publicKeyGen, err := x509util.GenerateKeyPublic(publicKeyCij)
					if err != nil {
						panic(err)
					}
					keyInfoItemPhase3 = &models.KeyInfoItemPhase3{
						DomainId:     user.DomainId,
						UserId:       user.UserId,
						PositionUser: user.Position,
						PositionItem: int32(j),
						KeyPrivate:   privateKeyGen,
						KeyPublic:    publicKeyGen,
						Part:         2,
					}
					err = inst.componentsContainer.KeyInfoItemPhase3Repository().Create(userContext, keyInfoItemPhase3)
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.Create").Error())
					}
				}
				privateKey, err := x509util.ExtractKeyPrivate(keyInfoItemPhase3.KeyPrivate)
				if err != nil {
					log.Printf(errutil.Wrap(err, "x509util.ExtractKeyPrivate").Error())
				}
				ite := int(math.Ceil(float64(j) / float64(combinedData.NumberItem2)))
				resultCard, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.DomainId,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.UserId,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    ite,
						},
					},
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "ResultCardRepository.FindLast").Error())
				}
				resultCardData := int32(0)
				if resultCard != nil {
					resultCardData = resultCard.PositionOption
				}
				itemByte := new(big.Int).SetInt64(int64(resultCardData))
				uiGx, uiGy := curve.ScalarBaseMult(itemByte.Bytes())
				// keyInfoXi
				keyInfoXi, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Value:    user.DomainId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Value:    user.UserId,
							Operator: entity.FindRequestFilterOperatorEqualTo,
						},
					},
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
				}
				publicKeyXi, err := x509util.ExtractKeyPublic(keyInfoXi.KeyPublic)
				if err != nil {
					panic(err)
				}
				valueprivateKey := new(big.Int).SetBytes(privateKey.D.Bytes())
				cXijx, cXijy := curve.ScalarMult(publicKeyXi.X, publicKeyXi.Y, valueprivateKey.Bytes())

				C1x, C1y := ellipticutil.AddPoint(curve, uiGx, uiGy, cXijx, cXijy)
				decryptC1 := elliptic.Marshal(curve, C1x, C1y)

				C2x, C2y := curve.ScalarBaseMult(valueprivateKey.Bytes())
				decryptC2 := elliptic.Marshal(curve, C2x, C2y)

				err = inst.componentsContainer.RecommendService().ProcessedDataSendPhase3TwoPart(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiPhase3TwoPart{
					UserId:               user.UserId,
					PositionUser:         user.Position,
					PositionItem:         int32(j),
					ProcessedDataC1:      hex.EncodeToString(decryptC1),
					ProcessedDataC2:      hex.EncodeToString(decryptC2),
					PositionItemOriginal: int32(ite),
					Part:                 2,
				})
				if err != nil {
					log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendPhase3TwoPart").Error())
				}
			}
			timeOneUserEnd := time.Now().UnixMilli()
			dentaTime := timeOneUserEnd - timeOneUserStart
			mutex.Lock()
			totalTimeSendDataPhase3 += dentaTime
			countUser++
			mutex.Unlock()
		}(user)
	}

	wg.Wait()
	fmt.Printf("Time send data phase 3 of %d user: %d \n", countUser, totalTimeSendDataPhase3)
	fmt.Printf("Time send data phase 3 of 1 user: %d \n", int(totalTimeSendDataPhase3)/countUser)
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) FakeSendProcessDataSurveyPhase4TwoPart(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	if combinedData.TenantId2 != req.Value {
		return nil, status.Errorf(codes.PermissionDenied, "Tenant is not allow create phase4 two part")
	}

	totalTimeSendData := int64(0)
	countUser := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			j := 1
			isStopj := false
			timeOneUserStart := time.Now().UnixMilli()
			processedDataSendPhase3, err := inst.componentsContainer.RecommendService().ProcessedDataSendPhase3Get(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiPhase3TwoPart{
				UserId:       user.UserId,
				PositionUser: user.Position,
				PositionItem: int32(j),
				Part:         2,
			})
			if err != nil {
				log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendPhase3Get").Error())
			}
			if processedDataSendPhase3 == nil {
				log.Printf(errutil.NewWithMessage("processedDataSendPhase3 Not found").Error())
			}
			pointC1, err := hex.DecodeString(processedDataSendPhase3.ProcessedDataC1)
			if err != nil {
				panic(err)
			}
			C1x, C1y := elliptic.Unmarshal(curve, pointC1)
			pointC2, err := hex.DecodeString(processedDataSendPhase3.ProcessedDataC2)
			if err != nil {
				panic(err)
			}
			C2x, C2y := elliptic.Unmarshal(curve, pointC2)

			for t := 1; t <= int(combinedData.SkTwoPart)-1; t++ {
				if isStopj {
					break
				}
				t2 := t % int(combinedData.NkTwoPart)
				if t%int(combinedData.NkTwoPart) == 0 {
					t2 = int(combinedData.NkTwoPart)
				}
				for k := t + 1; k <= int(combinedData.SkTwoPart); k++ {
					k2 := k % int(combinedData.NkTwoPart)
					if k%int(combinedData.NkTwoPart) == 0 {
						k2 = int(combinedData.NkTwoPart)
					}
					jte := int32(j) % combinedData.NumberItem2
					if int32(j)%combinedData.NumberItem2 == 0 {
						jte = combinedData.NumberItem2
					}
					resultCard, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    user.DomainId,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    user.UserId,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Operator: entity.FindRequestFilterOperatorEqualTo,
								Value:    jte,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "ResultCardRepository.FindLast").Error())
					}
					resultCardData := int32(0)
					if resultCard != nil {
						resultCardData = resultCard.PositionOption
					}
					itemByte := new(big.Int).SetInt64(int64(resultCardData))
					Cvijx, Cvijy := curve.ScalarMult(C1x, C1y, itemByte.Bytes())

					// keyInfoIK2
					keyInfoIK2, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    k2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}

					privateIK2, err := x509util.ExtractKeyPrivate(keyInfoIK2.KeyPrivate)
					if err != nil {
						panic(err)
					}

					//keyPublicT
					keyPublicT, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(t), 2)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}

					publicPt, err := x509util.ExtractKeyPublic(keyPublicT.KeyPublic)
					if err != nil {
						panic(err)
					}
					valuePrivateIK2 := new(big.Int).SetBytes(privateIK2.D.Bytes())
					ksuikKptx, ksuikKpty := curve.ScalarMult(publicPt.X, publicPt.Y, valuePrivateIK2.Bytes())

					R1x, R1y := ellipticutil.AddPoint(curve, Cvijx, Cvijy, ksuikKptx, ksuikKpty)

					privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
					if err != nil {
						panic(err)
					}

					// keyInfoYi
					keyInfoYi, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}
					privY, err := x509util.ExtractKeyPrivate(keyInfoYi.KeyPrivate)
					if err != nil {
						panic(err)
					}

					sharedPrivateKey := new(big.Int).Mul(privateKey.D, privY.D)

					C2rijx, C2rijy := curve.ScalarMult(C2x, C2y, sharedPrivateKey.Bytes())

					// keyInfoIT2
					keyInfoIT2, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    t2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}
					privateIT2, err := x509util.ExtractKeyPrivate(keyInfoIT2.KeyPrivate)
					if err != nil {
						panic(err)
					}

					//keyPublicK
					keyPublicK, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(k), 2)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}
					publicPk, err := x509util.ExtractKeyPublic(keyPublicK.KeyPublic)
					if err != nil {
						panic(err)
					}
					valuePrivateIT2 := new(big.Int).SetBytes(privateIT2.D.Bytes())
					ksuitKptx, ksuitKpty := curve.ScalarMult(publicPk.X, publicPk.Y, valuePrivateIT2.Bytes())

					R2x, R2y := ellipticutil.AddPoint(curve, C2rijx, C2rijy, ksuitKptx, ksuitKpty)

					publicKeyYi, err := x509util.ExtractKeyPublic(keyInfoYi.KeyPublic)
					if err != nil {
						panic(err)
					}
					valuePrivateij := new(big.Int).SetBytes(privateKey.D.Bytes())
					Yrijx, Yrijy := curve.ScalarMult(publicKeyYi.X, publicKeyYi.Y, valuePrivateij.Bytes())

					keyInfoXi, err := inst.componentsContainer.RecommendService().KeyPublicUserGet(userContext, combinedData.TenantId1, req.Ctx.TokenAgent, user.Position)
					publicKeyXi, err := x509util.ExtractKeyPublic(keyInfoXi.KeyPublic)
					if err != nil {
						panic(err)
					}
					Xvijx, Xvijy := curve.ScalarMult(publicKeyXi.X, publicKeyXi.Y, itemByte.Bytes())

					// invert one point
					XvijyY := new(big.Int).Neg(Xvijy)
					// point normalization
					XvijyYSub := new(big.Int).Mod(XvijyY, curve.Params().P)

					R3x, R3y := ellipticutil.AddPoint(curve, Yrijx, Yrijy, Xvijx, XvijyYSub)

					decryptR1 := elliptic.Marshal(curve, R1x, R1y)
					decryptR2 := elliptic.Marshal(curve, R2x, R2y)
					decryptR3 := elliptic.Marshal(curve, R3x, R3y)

					err = inst.componentsContainer.RecommendService().ProcessedDataSendPhase4TwoPart(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiPhase4TwoPart{
						UserId:                user.UserId,
						PositionUser:          user.Position,
						PositionItem:          int32(j),
						ProcessedDataR1:       hex.EncodeToString(decryptR1),
						ProcessedDataR2:       hex.EncodeToString(decryptR2),
						ProcessedDataR3:       hex.EncodeToString(decryptR3),
						PositionItemOriginal2: int32(jte),
						PositionItemOriginal1: processedDataSendPhase3.PositionItemOriginal,
						Part:                  2,
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendOnePart").Error())
					}

					if j == int(combinedData.STwoPart) {
						isStopj = true
						break
					} else {
						j++
					}
				}
			}
			timeOneUserEnd := time.Now().UnixMilli()
			dentaTime := timeOneUserEnd - timeOneUserStart
			mutex.Lock()
			totalTimeSendData += dentaTime
			countUser++
			mutex.Unlock()
		}(user)
	}
	wg.Wait()
	fmt.Printf("Time send data phase 4 of %d user: %d \n", countUser, totalTimeSendData)
	fmt.Printf("Time send data phase 4 of 1 user: %d \n", int(totalTimeSendData)/countUser)
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) FakeSendProcessDataSurveyPhase5TwoPart(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenant").Error())
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	if combinedData.TenantId1 != req.Value {
		return nil, status.Errorf(codes.PermissionDenied, "Tenant is not allow create phase5 two part")
	}
	totalTimeSendData := int64(0)
	countUser := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, user := range users {
		wg.Add(1)
		go func(user entity.User) {
			defer wg.Done()
			j := 1
			isStopj := false
			timeOneUserStart := time.Now().UnixMilli()
			processedDataSendPhase4, err := inst.componentsContainer.RecommendService().ProcessedDataSendPhase4Get(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiPhase4TwoPart{
				UserId:       user.UserId,
				PositionUser: user.Position,
				PositionItem: int32(j),
				Part:         2,
			})
			if err != nil {
				log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendPhase3Get").Error())
			}
			if processedDataSendPhase4 == nil {
				log.Printf(errutil.NewWithMessage("processedDataSendPhase4 Not found").Error())
			}
			pointR1, err := hex.DecodeString(processedDataSendPhase4.ProcessedDataR1)
			if err != nil {
				panic(err)
			}
			R1x, R1y := elliptic.Unmarshal(curve, pointR1)
			pointR2, err := hex.DecodeString(processedDataSendPhase4.ProcessedDataR2)
			if err != nil {
				panic(err)
			}
			R2x, R2y := elliptic.Unmarshal(curve, pointR2)
			pointR3, err := hex.DecodeString(processedDataSendPhase4.ProcessedDataR3)
			if err != nil {
				panic(err)
			}
			R3x, R3y := elliptic.Unmarshal(curve, pointR3)

			for t := 1; t <= int(combinedData.SkTwoPart)-1; t++ {
				if isStopj {
					break
				}
				t1 := int(math.Ceil(float64(t) / float64(int(combinedData.NkTwoPart))))
				for k := t + 1; k <= int(combinedData.SkTwoPart); k++ {
					k1 := int(math.Ceil(float64(k) / float64(int(combinedData.NkTwoPart))))

					keyInfoItemPhase3, err := inst.componentsContainer.KeyInfoItemPhase3Repository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    int32(j),
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}
					if keyInfoItemPhase3 == nil {
						log.Printf(errutil.NewWithMessage("keyInfoItemPhase3 not found").Error())
					}
					privateCij2, err := x509util.ExtractKeyPrivate(keyInfoItemPhase3.KeyPrivate)
					if err != nil {
						panic(err)
					}
					valuePrivateCij := new(big.Int).SetBytes(privateCij2.D.Bytes())
					cR3ijx, cR3ijy := curve.ScalarMult(R3x, R3y, valuePrivateCij.Bytes())

					R1cijR3x, R1cijR3y := ellipticutil.AddPoint(curve, R1x, R1y, cR3ijx, cR3ijy)

					// invert one point
					R2yY := new(big.Int).Neg(R2y)
					// point normalization
					R2yYSub := new(big.Int).Mod(R2yY, curve.Params().P)

					//3 cap dau R1 + cij*R3  - R2
					R1R3R2x, R1R3R2y := ellipticutil.AddPoint(curve, R1cijR3x, R1cijR3y, R2x, R2yYSub)

					// keyInfoIT1
					keyInfoIT1, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    t1,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}
					privateIT1, err := x509util.ExtractKeyPrivate(keyInfoIT1.KeyPrivate)
					if err != nil {
						panic(err)
					}

					//keyPublicK
					keyPublicK, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(k), 2)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}
					publicPk, err := x509util.ExtractKeyPublic(keyPublicK.KeyPublic)
					if err != nil {
						panic(err)
					}
					valuePrivateIT1 := new(big.Int).SetBytes(privateIT1.D.Bytes())
					ksuitKpkx, ksuitKpky := curve.ScalarMult(publicPk.X, publicPk.Y, valuePrivateIT1.Bytes())

					// invert one point
					ksuitKpkyY := new(big.Int).Neg(ksuitKpky)
					// point normalization
					ksuitKpkyYSub := new(big.Int).Mod(ksuitKpkyY, curve.Params().P)

					// 4 cap dau R1 + cij*R3  - R2 - ksuit1*KPk
					R1R3R2KPkx, R1R3R2KPky := ellipticutil.AddPoint(curve, R1R3R2x, R1R3R2y, ksuitKpkx, ksuitKpkyYSub)

					//keyInfoIK1
					keyInfoIK1, err := inst.componentsContainer.KeyInfoItemRepository().FindLast(userContext, &entity.FindRequest{
						Filters: []entity.FindRequestFilter{
							entity.FindRequestFilter{
								Key:      "DomainId",
								Value:    user.DomainId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "UserId",
								Value:    user.UserId,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "PositionItem",
								Value:    k1,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
							entity.FindRequestFilter{
								Key:      "Part",
								Value:    2,
								Operator: entity.FindRequestFilterOperatorEqualTo,
							},
						},
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "KeyInfoRepository.FindLast").Error())
					}
					privateIK1, err := x509util.ExtractKeyPrivate(keyInfoIK1.KeyPrivate)
					if err != nil {
						panic(err)
					}

					//keyPublicT
					keyPublicT, err := inst.componentsContainer.RecommendService().KeyPublicUse(userContext, user.DomainId, req.Ctx.TokenAgent, int32(t), 2)
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.KeyPublicUse").Error())
					}
					publicPt, err := x509util.ExtractKeyPublic(keyPublicT.KeyPublic)
					if err != nil {
						panic(err)
					}
					valuePrivateIK1 := new(big.Int).SetBytes(privateIK1.D.Bytes())
					ksuikKptx, ksuikKpty := curve.ScalarMult(publicPt.X, publicPt.Y, valuePrivateIK1.Bytes())

					Mijx, MijY := ellipticutil.AddPoint(curve, R1R3R2KPkx, R1R3R2KPky, ksuikKptx, ksuikKpty)

					decrypt := elliptic.Marshal(curve, Mijx, MijY)

					err = inst.componentsContainer.RecommendService().ProcessedDataSendPart(userContext, user.DomainId, req.Ctx.TokenAgent, &recommendApi.RecommendApiResultCard{
						UserId:                user.UserId,
						PositionUser:          user.Position,
						PositionItem:          int32(j),
						ProcessedData:         hex.EncodeToString(decrypt),
						PositionItemOriginal1: processedDataSendPhase4.PositionItemOriginal1,
						PositionItemOriginal2: processedDataSendPhase4.PositionItemOriginal2,
						Part:                  2,
					})
					if err != nil {
						log.Printf(errutil.Wrap(err, "RecommendService.ProcessedDataSendPart").Error())
					}
					if j == int(combinedData.STwoPart) {
						isStopj = true
						break
					} else {
						j++
					}
				}
			}
			timeOneUserEnd := time.Now().UnixMilli()
			dentaTime := timeOneUserEnd - timeOneUserStart
			mutex.Lock()
			totalTimeSendData += dentaTime
			countUser++
			mutex.Unlock()
		}(user)
	}
	wg.Wait()
	fmt.Printf("Time send data phase 5 of %d user: %d \n", countUser, totalTimeSendData)
	fmt.Printf("Time send data phase 5 of 1 user: %d \n", int(totalTimeSendData)/countUser)
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}

func (inst *AgentpcService) ProcessDataSurvey2(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "CustomerService.GetTenantById").Error())
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		log.Printf(errutil.Wrap(err, "RecommendService.GetCombinedData").Error())
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}

	for j := 1; j <= int(combinedData.STwoPart); j++ {
		sumRating := int32(0)
		positionOriginal1 := int32(0)
		positionOriginal2 := int32(0)
		for i := 1; i <= int(combinedData.NumberUser); i++ {
			ju := int(math.Ceil(float64(j) / float64(combinedData.NumberItem2)))
			resultCard1, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
				Filters: []entity.FindRequestFilter{
					entity.FindRequestFilter{
						Key:      "DomainId",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    combinedData.TenantId1,
					},
					entity.FindRequestFilter{
						Key:      "PositionUser",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    i,
					},
					entity.FindRequestFilter{
						Key:      "PositionItem",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    ju,
					},
				},
			})
			if err != nil {
				log.Printf(errutil.Wrap(err, "ResultCardRepository.FindLast").Error())
			}
			jv := int32(j) % combinedData.NumberItem2
			if int32(j)%combinedData.NumberItem2 == 0 {
				jv = combinedData.NumberItem2
			}
			resultCard2, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
				Filters: []entity.FindRequestFilter{
					entity.FindRequestFilter{
						Key:      "DomainId",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    combinedData.TenantId2,
					},
					entity.FindRequestFilter{
						Key:      "PositionUser",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    i,
					},
					entity.FindRequestFilter{
						Key:      "PositionItem",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    jv,
					},
				},
			})
			if err != nil {
				log.Printf(errutil.Wrap(err, "ResultCardRepository.FindLast").Error())
			}
			sumRating += resultCard2.PositionOption * resultCard1.PositionOption
			positionOriginal1 = int32(ju)
			positionOriginal2 = int32(jv)
		}

		err = inst.componentsContainer.RecommendService().ProcessedDataSend2(userContext, req.Value, req.Ctx.TokenAgent, &recommendApi.RecommendApiProcessDataSurvey2{
			PositionItem:          int32(j),
			ProcessedData:         sumRating,
			PositionItemOriginal1: positionOriginal1,
			PositionItemOriginal2: positionOriginal2,
			Part:                  2,
		})
	}
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}
