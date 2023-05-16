package util

import (
	"crypto/rand"
	"encoding/json"
	"strconv"

	"github.com/thteam47/go-agent-pc/errutil"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func FromMessage(from protoreflect.ProtoMessage, to interface{}) error {
	marshaller := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	data, err := marshaller.Marshal(from)
	if err != nil {
		return errutil.Wrapf(err, "marshaller.Marshal")
	}
	err = json.Unmarshal(data, to)
	if err != nil {
		return errutil.Wrapf(err, "json.Unmarshal")
	}
	return nil
}

func ToMessage(from interface{}, to protoreflect.ProtoMessage) error {
	data, err := json.Marshal(from)
	if err != nil {
		return errutil.Wrapf(err, "json.Marshal")
	}
	unmarshaller := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err = unmarshaller.Unmarshal(data, to)
	if err != nil {
		return errutil.Wrapf(err, "unmarshaller.Unmarshal")
	}
	return nil
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func GenerateCodeOtp() (int, error) {
	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return 0, errutil.Wrapf(err, "rand.Read")
	}

	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + (codes[i] % 10))
	}

	return strconv.Atoi(string(codes))
}
