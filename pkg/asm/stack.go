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
package asm

type stack struct {
	vals []bool
}

func (s *stack) push(v bool) {
	s.vals = append(s.vals, v)
}

func (s *stack) len() int {
	return len(s.vals)
}

func (s *stack) top() bool {
	return s.vals[len(s.vals)-1]
}

func (s *stack) pop() bool {
	res := s.vals[len(s.vals)-1]
	s.vals = s.vals[0 : len(s.vals)-1]
	return res
}
