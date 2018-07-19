// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
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


package dotfiles

// func collectErrors(errs []error) error {
// 	if len(errs) == 0 {
// 		return nil
// 	}

// 	var errMsg string

// 	for _, err := range errs {
// 		errMsg += err.Error() + "\n"
// 	}

// 	return errors.New(errMsg)
// }

// func (p Profile) Init() error {
// 	for _, location := range p.Locations {
// 		err := os.Mkdir(location, os.ModePerm)
// 		if err != nil {
// 			return err
// 		}

// 		backend := loadBackend(p.Backend)
// 		err = backend.NewProfile(location)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
