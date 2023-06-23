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

// RevocationPluginConfigInfo revocation plugin config info
//
// swagger:model RevocationPluginConfigInfo
type RevocationPluginConfigInfo struct {

	// app config
	AppConfig *AppConfig `json:"appConfig,omitempty"`

	// custom config
	CustomConfig *PublicCustomConfigInfo `json:"customConfig,omitempty"`

	// extend type
	// Enum: [APP CUSTOM]
	ExtendType string `json:"extendType,omitempty"`

	// namespace
	// Required: true
	Namespace *string `json:"namespace"`
}

// Validate validates this revocation plugin config info
func (m *RevocationPluginConfigInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAppConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCustomConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExtendType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNamespace(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RevocationPluginConfigInfo) validateAppConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.AppConfig) { // not required
		return nil
	}

	if m.AppConfig != nil {
		if err := m.AppConfig.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("appConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("appConfig")
			}
			return err
		}
	}

	return nil
}

func (m *RevocationPluginConfigInfo) validateCustomConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.CustomConfig) { // not required
		return nil
	}

	if m.CustomConfig != nil {
		if err := m.CustomConfig.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("customConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("customConfig")
			}
			return err
		}
	}

	return nil
}

var revocationPluginConfigInfoTypeExtendTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["APP","CUSTOM"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		revocationPluginConfigInfoTypeExtendTypePropEnum = append(revocationPluginConfigInfoTypeExtendTypePropEnum, v)
	}
}

const (

	// RevocationPluginConfigInfoExtendTypeAPP captures enum value "APP"
	RevocationPluginConfigInfoExtendTypeAPP string = "APP"

	// RevocationPluginConfigInfoExtendTypeCUSTOM captures enum value "CUSTOM"
	RevocationPluginConfigInfoExtendTypeCUSTOM string = "CUSTOM"
)

// prop value enum
func (m *RevocationPluginConfigInfo) validateExtendTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, revocationPluginConfigInfoTypeExtendTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *RevocationPluginConfigInfo) validateExtendType(formats strfmt.Registry) error {
	if swag.IsZero(m.ExtendType) { // not required
		return nil
	}

	// value enum
	if err := m.validateExtendTypeEnum("extendType", "body", m.ExtendType); err != nil {
		return err
	}

	return nil
}

func (m *RevocationPluginConfigInfo) validateNamespace(formats strfmt.Registry) error {

	if err := validate.Required("namespace", "body", m.Namespace); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this revocation plugin config info based on the context it is used
func (m *RevocationPluginConfigInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAppConfig(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCustomConfig(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RevocationPluginConfigInfo) contextValidateAppConfig(ctx context.Context, formats strfmt.Registry) error {

	if m.AppConfig != nil {
		if err := m.AppConfig.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("appConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("appConfig")
			}
			return err
		}
	}

	return nil
}

func (m *RevocationPluginConfigInfo) contextValidateCustomConfig(ctx context.Context, formats strfmt.Registry) error {

	if m.CustomConfig != nil {
		if err := m.CustomConfig.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("customConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("customConfig")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RevocationPluginConfigInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RevocationPluginConfigInfo) UnmarshalBinary(b []byte) error {
	var res RevocationPluginConfigInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}