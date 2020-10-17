/*
 * Copyright (c) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of cbmasm.
 *
 * cbmasm is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * cbmasm is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.
 */
package errors

import (
	"fmt"

	"github.com/asig/cbmasm/pkg/text"
)

type Error struct {
	Pos text.Pos
	Msg string
}

func (e Error) String() string {
	return fmt.Sprintf("%s, line %d, col %d: %s", e.Pos.Filename, e.Pos.Line, e.Pos.Col, e.Msg)
}

type Sink interface {
	AddError(pos text.Pos, message string, args ...interface{})
}

type Modifier interface {
	Modify(err Error) Error
}
