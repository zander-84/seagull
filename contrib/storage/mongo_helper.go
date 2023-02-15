package storage

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MongoPkAny2ObjectIDs(in []any) ([]primitive.ObjectID, error) {
	out := make([]primitive.ObjectID, 0, len(in))
	for _, v := range in {
		id, err := MongoPkAny2ObjectID(v)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

func MongoPkAny2ObjectID(in any) (primitive.ObjectID, error) {
	in2, ok := in.(string)
	if !ok {
		return primitive.ObjectID{}, errors.New("id must a number")
	}
	return MongoPkString2ObjectID(in2)
}

func MongoPkString2ObjectID(in string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(in)
}
func MongoPkObjectID2String(in primitive.ObjectID) string {
	return in.Hex()
}
