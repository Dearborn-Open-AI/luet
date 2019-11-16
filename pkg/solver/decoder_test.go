// Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package solver_test

import (
	"strconv"

	pkg "github.com/mudler/luet/pkg/package"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mudler/luet/pkg/solver"
)

var _ = Describe("Decoder", func() {
	Context("Assertion ordering", func() {
		eq := 0
		for index := 0; index < 300; index++ { // Just to make sure we don't have false positives
			It("Orders them correctly #"+strconv.Itoa(index), func() {
				db := pkg.NewInMemoryDatabase(false)

				C := pkg.NewPackage("C", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				E := pkg.NewPackage("E", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				F := pkg.NewPackage("F", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				G := pkg.NewPackage("G", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				H := pkg.NewPackage("H", "", []*pkg.DefaultPackage{G}, []*pkg.DefaultPackage{})
				D := pkg.NewPackage("D", "", []*pkg.DefaultPackage{H}, []*pkg.DefaultPackage{})
				B := pkg.NewPackage("B", "", []*pkg.DefaultPackage{D}, []*pkg.DefaultPackage{})
				A := pkg.NewPackage("A", "", []*pkg.DefaultPackage{B}, []*pkg.DefaultPackage{})

				s := NewSolver([]pkg.Package{C}, []pkg.Package{A, B, C, D, E, F, G}, db)

				solution, err := s.Install([]pkg.Package{A})
				Expect(solution).To(ContainElement(PackageAssert{Package: A.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: B.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: D.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: C.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: H.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: G.IsFlagged(true), Value: true}))

				Expect(len(solution)).To(Equal(6))
				Expect(err).ToNot(HaveOccurred())
				solution = solution.Order(A.GetFingerPrint())
				//	Expect(len(solution)).To(Equal(6))
				Expect(solution[0].Package.GetName()).To(Equal("G"))
				Expect(solution[1].Package.GetName()).To(Equal("H"))
				Expect(solution[2].Package.GetName()).To(Equal("D"))
				Expect(solution[3].Package.GetName()).To(Equal("B"))
				eq += len(solution)
				//Expect(solution[4].Package.GetName()).To(Equal("A"))
				//Expect(solution[5].Package.GetName()).To(Equal("C")) As C doesn't have any dep it can be in both positions
			})
		}

		It("Expects perfect equality when ordered", func() {
			Expect(eq).To(Equal(300 * 5)) // assertions lenghts
		})

		disequality := 0
		equality := 0
		for index := 0; index < 300; index++ { // Just to make sure we don't have false positives
			It("Doesn't order them correctly otherwise #"+strconv.Itoa(index), func() {
				db := pkg.NewInMemoryDatabase(false)

				C := pkg.NewPackage("C", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				E := pkg.NewPackage("E", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				F := pkg.NewPackage("F", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				G := pkg.NewPackage("G", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
				H := pkg.NewPackage("H", "", []*pkg.DefaultPackage{G}, []*pkg.DefaultPackage{})
				D := pkg.NewPackage("D", "", []*pkg.DefaultPackage{H}, []*pkg.DefaultPackage{})
				B := pkg.NewPackage("B", "", []*pkg.DefaultPackage{D}, []*pkg.DefaultPackage{})
				A := pkg.NewPackage("A", "", []*pkg.DefaultPackage{B}, []*pkg.DefaultPackage{})

				s := NewSolver([]pkg.Package{C}, []pkg.Package{A, B, C, D, E, F, G}, db)

				solution, err := s.Install([]pkg.Package{A})
				Expect(solution).To(ContainElement(PackageAssert{Package: A.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: B.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: D.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: C.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: H.IsFlagged(true), Value: true}))
				Expect(solution).To(ContainElement(PackageAssert{Package: G.IsFlagged(true), Value: true}))

				Expect(len(solution)).To(Equal(6))
				Expect(err).ToNot(HaveOccurred())
				if solution[0].Package.GetName() != "G" {
					disequality++
				} else {
					equality++
				}
				if solution[1].Package.GetName() != "H" {
					disequality++
				} else {
					equality++
				}
				if solution[2].Package.GetName() != "D" {
					disequality++
				} else {
					equality++
				}
				if solution[3].Package.GetName() != "B" {
					disequality++
				} else {
					equality++
				}
				if solution[4].Package.GetName() != "A" {
					disequality++
				} else {
					equality++
				}

			})
			It("Expect disequality", func() {
				Expect(disequality).ToNot(Equal(0))
				Expect(equality).ToNot(Equal(300 * 6))
			})
		}
	})

	Context("Assertion hashing", func() {
		It("Hashes them, and could be used for comparison", func() {
			db := pkg.NewInMemoryDatabase(false)

			C := pkg.NewPackage("C", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
			E := pkg.NewPackage("E", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
			F := pkg.NewPackage("F", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
			G := pkg.NewPackage("G", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
			H := pkg.NewPackage("H", "", []*pkg.DefaultPackage{G}, []*pkg.DefaultPackage{})
			D := pkg.NewPackage("D", "", []*pkg.DefaultPackage{H}, []*pkg.DefaultPackage{})
			B := pkg.NewPackage("B", "", []*pkg.DefaultPackage{D}, []*pkg.DefaultPackage{})
			A := pkg.NewPackage("A", "", []*pkg.DefaultPackage{B}, []*pkg.DefaultPackage{})

			s := NewSolver([]pkg.Package{C}, []pkg.Package{A, B, C, D, E, F, G}, db)

			solution, err := s.Install([]pkg.Package{A})
			Expect(solution).To(ContainElement(PackageAssert{Package: A.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: B.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: D.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: C.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: H.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: G.IsFlagged(true), Value: true}))

			Expect(len(solution)).To(Equal(6))
			Expect(err).ToNot(HaveOccurred())
			solution = solution.Order(A.GetFingerPrint())
			//	Expect(len(solution)).To(Equal(6))
			Expect(solution[0].Package.GetName()).To(Equal("G"))
			Expect(solution[1].Package.GetName()).To(Equal("H"))
			Expect(solution[2].Package.GetName()).To(Equal("D"))
			Expect(solution[3].Package.GetName()).To(Equal("B"))

			hash := solution.AssertionHash()

			solution, err = s.Install([]pkg.Package{B})
			Expect(solution).To(ContainElement(PackageAssert{Package: B.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: D.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: C.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: H.IsFlagged(true), Value: true}))
			Expect(solution).To(ContainElement(PackageAssert{Package: G.IsFlagged(true), Value: true}))

			Expect(len(solution)).To(Equal(6))
			Expect(err).ToNot(HaveOccurred())
			solution = solution.Order(B.GetFingerPrint())
			hash2 := solution.AssertionHash()

			//	Expect(len(solution)).To(Equal(6))
			Expect(solution[0].Package.GetName()).To(Equal("A"))
			Expect(solution[1].Package.GetName()).To(Equal("G"))
			Expect(solution[2].Package.GetName()).To(Equal("H"))
			Expect(solution[3].Package.GetName()).To(Equal("D"))
			Expect(solution[4].Package.GetName()).To(Equal("B"))
			Expect(solution[0].Value).ToNot(BeTrue())
			Expect(solution[0].Package.Flagged()).To(BeTrue())

			Expect(hash).ToNot(Equal(""))
			Expect(hash2).ToNot(Equal(""))
			Expect(hash != hash2).To(BeTrue())
			db2 := pkg.NewInMemoryDatabase(false)

			X := pkg.NewPackage("X", "", []*pkg.DefaultPackage{}, []*pkg.DefaultPackage{})
			Y := pkg.NewPackage("Y", "", []*pkg.DefaultPackage{X}, []*pkg.DefaultPackage{})
			Z := pkg.NewPackage("Z", "", []*pkg.DefaultPackage{X}, []*pkg.DefaultPackage{})
			s = NewSolver([]pkg.Package{}, []pkg.Package{X, Y, Z}, db2)
			solution, err = s.Install([]pkg.Package{Y})
			Expect(err).ToNot(HaveOccurred())

			solution2, err := s.Install([]pkg.Package{Z})
			Expect(err).ToNot(HaveOccurred())
			Expect(solution.Order(Y.GetFingerPrint()).Drop(Y).AssertionHash() == solution2.Order(Z.GetFingerPrint()).Drop(Z).AssertionHash()).To(BeTrue())
		})

	})
})