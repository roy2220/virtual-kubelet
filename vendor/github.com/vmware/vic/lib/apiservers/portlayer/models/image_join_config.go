package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// ImageJoinConfig image join config
// swagger:model ImageJoinConfig
type ImageJoinConfig struct {

	// delta ID
	// Required: true
	DeltaID string `json:"deltaID"`

	// handle
	// Required: true
	Handle interface{} `json:"handle"`

	// image ID
	ImageID string `json:"imageID,omitempty"`

	// repo name
	RepoName string `json:"repoName,omitempty"`
}

// Validate validates this image join config
func (m *ImageJoinConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDeltaID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateHandle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImageJoinConfig) validateDeltaID(formats strfmt.Registry) error {

	if err := validate.RequiredString("deltaID", "body", string(m.DeltaID)); err != nil {
		return err
	}

	return nil
}

func (m *ImageJoinConfig) validateHandle(formats strfmt.Registry) error {

	return nil
}