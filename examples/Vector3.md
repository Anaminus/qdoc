# Vector3
[Vector3]: #user-content-vector3
```
type Vector3
```

Vector3 describes a vector in 3D space.

<table>
<thead><tr><th>Table of Contents</th></tr></thead>
<tbody><tr><td>

1. [Vector3][Vector3]
	1. [Vector3.new][Vector3.new]
	2. [Vector3.X][Vector3.X]
	3. [Vector3.Y][Vector3.Y]
	4. [Vector3.Z][Vector3.Z]
	5. [Vector3.Cross][Vector3.Cross]
	6. [Vector3.Dot][Vector3.Dot]
	7. [Vector3.Lerp][Vector3.Lerp]
	8. [Vector3.Magnitude][Vector3.Magnitude]
	9. [Vector3.Unit][Vector3.Unit]
	10. [Addition][Vector3.__add]
	11. [Division][Vector3.__div]
	12. [Equality][Vector3.__eq]
	13. [Multiplication][Vector3.__mul]
	14. [Subtraction][Vector3.__sub]
	15. [String][Vector3.__tostring]
	16. [Negation][Vector3.__unm]

</td></tr></tbody>
</table>

## Vector3.new
[Vector3.new]: #user-content-vector3new
```
Vector3.new(x: number?, y: number?, z: number?): Vector3
```

new returns a new Vector3 with the given components.

## Vector3.X
[Vector3.X]: #user-content-vector3x
```
Vector3.X: number
```

X returns the x component of the vector.

## Vector3.Y
[Vector3.Y]: #user-content-vector3y
```
Vector3.Y: number
```

Y returns the y component of the vector.

## Vector3.Z
[Vector3.Z]: #user-content-vector3z
```
Vector3.Z: number
```

Z returns the z component of the vector.

## Vector3.Cross
[Vector3.Cross]: #user-content-vector3cross
```
Vector3:Cross(v: Vector3): Vector3
```

Cross returns the cross product of the two vectors.

## Vector3.Dot
[Vector3.Dot]: #user-content-vector3dot
```
Vector3:Dot(v: Vector3): number
```

Dot returns the scalar dot product of the two vectors.

## Vector3.Lerp
[Vector3.Lerp]: #user-content-vector3lerp
```
Vector3:Lerp(goal: Vector3, alpha: number): Vector3
```

Lerp returns a vector linearly interpolated between the vector and
*goal* by the fraction *alpha*.

## Vector3.Magnitude
[Vector3.Magnitude]: #user-content-vector3magnitude
```
Vector3:Magnitude(): number
```

Magnitude returns the length of the vector.

## Vector3.Unit
[Vector3.Unit]: #user-content-vector3unit
```
Vector3:Unit(): Vector3
```

Unit returns a vector with the same direction, but with a length of 1.

## Addition
[Vector3.__add]: #user-content-addition
```
Vector3 + Vector3: Vector3
```

Adding two vectors returns a Vector3 with the sum of each of their
components.

## Division
[Vector3.__div]: #user-content-division
```
Vector3 / Vector3: Vector3
Vector3 / number: Vector3
number / Vector3: Vector3
```

Dividing two vectors returns a Vector3 with each component of the first
divided by those of the second.

A vector may also be divided by a number, which returns a Vector3 with each
component divided by the number.

## Equality
[Vector3.__eq]: #user-content-equality
```
Vector3 == Vector3: boolean
```

Two vectors are equal if each of their [X][Vector3.X], [Y][Vector3.Y],
and [Z][Vector3.Z] components are equal.

## Multiplication
[Vector3.__mul]: #user-content-multiplication
```
Vector3 * Vector3: Vector3
Vector3 * number: Vector3
number * Vector3: Vector3
```

Multiplying two vectors returns a Vector3 with each component of the
first multipied by those of the second.

A vector may also be multiplied by a number, which returns a Vector3 with
each component multiplied by the number.

## Subtraction
[Vector3.__sub]: #user-content-subtraction
```
Vector3 - Vector3: Vector3
```

Subtracting two vectors returns a Vector3 with each component of the
second subtracted from those of the first.

## String
[Vector3.__tostring]: #user-content-string

A vector, when converted to a string, displays each of its
[X][Vector3.X], [Y][Vector3.Y], and [Z][Vector3.Z] components.

## Negation
[Vector3.__unm]: #user-content-negation
```
-Vector3: Vector3
```

Negating a vector returns a Vector3 with each the negation of each of
component.

