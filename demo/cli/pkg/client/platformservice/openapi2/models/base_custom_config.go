// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BaseCustomConfig base custom config
//
// swagger:model BaseCustomConfig
type BaseCustomConfig struct {

	// connect type: INSECURE, TLS, default is INSECURE
	// Required: true
	// Enum: [INSECURE TLS]
	ConnectionType *string `json:"connectionType"`

	// plugin grpc server address: <host>:<port>
	// Required: true
	GrpcServerAddress *string `json:"grpcServerAddress"`
}

// Validate validates this base custom config
func (m *BaseCustomConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConnectionType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGrpcServerAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var baseCustomConfigTypeConnectionTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["INSECURE","TLS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		baseCustomConfigTypeConnectionTypePropEnum = append(baseCustomConfigTypeConnectionTypePropEnum, v)
	}
}

const (

	// BaseCustomConfigConnectionTypeINSECURE captures enum value "INSECURE"
	BaseCustomConfigConnectionTypeINSECURE string = "INSECURE"

	// BaseCustomConfigConnectionTypeTLS captures enum value "TLS"
	BaseCustomConfigConnectionTypeTLS string = "TLS"
)

// prop value enum
func (m *BaseCustomConfig) validateConnectionTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, baseCustomConfigTypeConnectionTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *BaseCustomConfig) validateConnectionType(formats strfmt.Registry) error {

	if err := validate.Required("connectionType", "body", m.ConnectionType); err != nil {
		return err
	}

	// value enum
	if err := m.validateConnectionTypeEnum("connectionType", "body", *m.ConnectionType); err != nil {
		return err
	}

	return nil
}

func (m *BaseCustomConfig) validateGrpcServerAddress(formats strfmt.Registry) error {

	if err := validate.Required("grpcServerAddress", "body", m.GrpcServerAddress); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this base custom config based on context it is used
func (m *BaseCustomConfig) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *BaseCustomConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BaseCustomConfig) UnmarshalBinary(b []byte) error {
	var res BaseCustomConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}