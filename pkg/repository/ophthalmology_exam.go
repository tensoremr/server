/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package repository

import (
	"github.com/tensoremr/server/pkg/models"
	"gorm.io/gorm"
)

type OpthalmologyExamRepository struct {
	DB *gorm.DB
}

func ProvideOpthalmologyExamRepository(DB *gorm.DB) OpthalmologyExamRepository {
	return OpthalmologyExamRepository{DB: DB}
}


// Save ...
func (r *OpthalmologyExamRepository) Save(m *models.OpthalmologyExam) error {
	return r.DB.Create(&m).Error
}

// Get ...
func (r *OpthalmologyExamRepository) Get(m *models.OpthalmologyExam, filter models.OpthalmologyExam) error {
	return r.DB.Where(filter).Take(&m).Error
}

// GetByPatientChart ...
func (r *OpthalmologyExamRepository) GetByPatientChart(m *models.OpthalmologyExam, ID int) error {
	return r.DB.Where("patient_chart_id = ?", ID).Take(&m).Error
}

// Update ...
func (r *OpthalmologyExamRepository) Update(m *models.OpthalmologyExam) error {
	return r.DB.Updates(&m).Error
}
