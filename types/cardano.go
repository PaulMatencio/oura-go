package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

type Cardano uint64

func CardanoType() *bsoncodec.Registry {
	// create a custom registry builder
	rb := bsoncodec.NewRegistryBuilder()

	// register default codecs and encoders/decoders
	var primitiveCodecs bson.PrimitiveCodecs
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	primitiveCodecs.RegisterPrimitiveCodecs(rb)

	// register custom encoder/decoder
	myCardanoType := reflect.TypeOf(Cardano(0))

	rb.RegisterTypeEncoder(
		myCardanoType,
		bsoncodec.ValueEncoderFunc(func(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
			if !val.IsValid() || val.Type() != myCardanoType {
				return bsoncodec.ValueEncoderError{
					Name:     "MyCardanoEncodeValue",
					Types:    []reflect.Type{myCardanoType},
					Received: val,
				}
			}
			// IMPORTANT STEP: cast uint64 to int64 so it can be stored in mongo
			vw.WriteInt64(int64(val.Uint()))
			return nil
		}),
	)

	rb.RegisterTypeDecoder(
		myCardanoType,
		bsoncodec.ValueDecoderFunc(func(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
			// IMPORTANT STEP: read store value in mongo as int64
			read, err := vr.ReadInt64()
			if err != nil {
				return err
			}
			// IMPORTANT STEP: cast back to uint64
			val.SetUint(uint64(read))
			return nil
		}),
	)

	// build the registry
	return rb.Build()
}
