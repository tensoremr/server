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

type PastOptSurgeryRepository struct {
	DB *gorm.DB
}

func ProvidePastOptSurgeryRepository(DB *gorm.DB) PastOptSurgeryRepository {
	return PastOptSurgeryRepository{DB: DB}
}

// Save ...
func (r *PastOptSurgeryRepository) Save(m *models.PastOptSurgery) error {
	return r.DB.Create(&m).Error
}

// Get ...
func (r *PastOptSurgeryRepository) Get(m *models.PastOptSurgery, ID int) error {
	return r.DB.Where("id = ?", ID).Take(&m).Error
}

// Update ...
func (r *PastOptSurgeryRepository) Update(m *models.PastOptSurgery) error {
	return r.DB.Save(&m).Error
}

// Delete ...
func (r *PastOptSurgeryRepository) Delete(ID int) error {
	return r.DB.Where("id = ?", ID).Delete(&models.PastOptSurgery{}).Error
}
