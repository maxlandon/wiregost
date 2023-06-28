package db

// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"gorm.io/gorm"

	"github.com/maxlandon/wiregost/internal/server/db/models"
)

// ErrRecordNotFound - Record not found error
var ErrRecordNotFound = gorm.ErrRecordNotFound

// OperatorByToken - Select an operator by token value
func OperatorByToken(value string) (*models.Operator, error) {
	if len(value) < 1 {
		return nil, ErrRecordNotFound
	}
	operator := &models.Operator{}
	err := Session().Where(&models.Operator{
		Token: value,
	}).First(operator).Error
	return operator, err
}

// OperatorAll - Select all operators from the database
func OperatorAll() ([]*models.Operator, error) {
	operators := []*models.Operator{}
	err := Session().Distinct("Name").Find(&operators).Error
	return operators, err
}

// GetKeyValue - Get a value from a key
func GetKeyValue(key string) (string, error) {
	keyValue := &models.KeyValue{}
	err := Session().Where(&models.KeyValue{
		Key: key,
	}).First(keyValue).Error
	return keyValue.Value, err
}

// SetKeyValue - Set the value for a key/value pair
func SetKeyValue(key string, value string) error {
	err := Session().Where(&models.KeyValue{
		Key: key,
	}).First(&models.KeyValue{}).Error
	if err == ErrRecordNotFound {
		err = Session().Create(&models.KeyValue{
			Key:   key,
			Value: value,
		}).Error
	} else {
		err = Session().Where(&models.KeyValue{
			Key: key,
		}).Updates(models.KeyValue{
			Key:   key,
			Value: value,
		}).Error
	}
	return err
}

// DeleteKeyValue - Delete a key/value pair
func DeleteKeyValue(key string, value string) error {
	return Session().Delete(&models.KeyValue{
		Key: key,
	}).Error
}
