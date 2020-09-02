--@sec: Vector3
--@def: type Vector3
--@doc: Vector3 describes a vector in 3D space.
local Vector3 = {__index={}}

local function newVector3(x, y, z)
	return setmetatable({X=x, Y=y, Z=z}, Vector3)
end

--@sec: Vector3.X
--@ord: -1
--@def: Vector3.X: number
--@doc: X returns the x component of the vector.

--@sec: Vector3.Y
--@ord: -1
--@def: Vector3.Y: number
--@doc: Y returns the y component of the vector.

--@sec: Vector3.Z
--@ord: -1
--@def: Vector3.Z: number
--@doc: Z returns the z component of the vector.

--@sec: Vector3.Magnitude
--@def: Vector3:Magnitude(): number
--@doc: Magnitude returns the length of the vector.
function Vector3.__index:Magnitude()
	return math.sqrt(self.X*self.X + self.Y*self.Y + self.Z*self.Z)
end

--@sec: Vector3.Unit
--@def: Vector3:Unit(): Vector3
--@doc: Unit returns a vector with the same direction, but with a length of 1.
function Vector3.__index:Unit()
	local m = self:Magnitude()
	return newVector3(self.X/m, self.Y/m, self.Z/m)
end

--@sec: Vector3.Lerp
--@def: Vector3:Lerp(goal: Vector3, alpha: number): Vector3
--@doc: Lerp returns a vector linearly interpolated between the vector and
-- *goal* by the fraction *alpha*.
function Vector3.__index:Lerp(v, alpha)
	return (1-alpha)*self + alpha*v
end

--@sec: Vector3.Dot
--@def: Vector3:Dot(v: Vector3): number
--@doc: Dot returns the scalar dot product of the two vectors.
function Vector3.__index:Dot(v)
	return self.X*v.X + self.Y*v.Y + self.Z*v.Z
end

--@sec: Vector3.Cross
--@def: Vector3:Cross(v: Vector3): Vector3
--@doc: Cross returns the cross product of the two vectors.
function Vector3.__index:Cross(v)
	return newVector3(
		self.y * v.z - self.z * v.y,
		self.z * v.x - self.x * v.z,
		self.x * v.y - self.y * v.x
	)
end

function Vector3:__newindex(field)
	error(field .. " cannot be assigned to", 2)
end

--@sec: Equality {Vector3.__eq}
--@ord: 1
--@def: Vector3 == Vector3: boolean
--@doc: Two vectors are equal if each of their X, Y, and Z components are equal.
function Vector3:__eq(v)
	return self.X == v.X and self.Y == v.Y and self.Z == v.Z
end

--@sec: String {Vector3.__tostring}
--@ord: 1
--@doc: A vector, when converted to a string, displays each of its X, Y, and Z
-- components.
function Vector3:__tostring()
	return self.X .. ", " .. self.Y .. ", " ..self.Z
end

--@sec: Addition {Vector3.__add}
--@ord: 1
--@def: Vector3 + Vector3: Vector3
--@doc: Adding two vectors returns a Vector3 with the sum of each of their
-- components.
function Vector3:__add(v)
	return newVector3(
		self.X + v.X,
		self.Y + v.Y,
		self.Z + v.Z
	)
end

--@sec: Subtraction {Vector3.__sub}
--@ord: 1
--@def: Vector3 - Vector3: Vector3
--@doc: Subtracting two vectors returns a Vector3 with each component of the
-- second subtracted from those of the first.
function Vector3:__sub(v)
	return newVector3(
		self.X - v.X,
		self.Y - v.Y,
		self.Z - v.Z
	)
end

--@sec: Multiplication {Vector3.__mul}
--@ord: 1
--@def:
-- Vector3 * Vector3: Vector3
-- Vector3 * number: Vector3
-- number * Vector3: Vector3
--@doc: Multiplying two vectors returns a Vector3 with each component of the
-- first multipied by those of the second.
--
-- A vector may also be multiplied by a number, which returns a Vector3 with
-- each component multiplied by the number.
function Vector3.__mul(a, b)
	if type(a) == "number" then
		return newVector3(a*b.X, a*b.Y, a*b.Z)
	elseif type(b) == "number" then
		return newVector3(a.X*b, a.Y*b, a.Z*b)
	end
	return newVector3(a.X*b.X, a.Y*b.Y, a.Z*b.Z)
end

--@sec: Division {Vector3.__div}
--@ord: 1
--@def:
-- Vector3 / Vector3: Vector3
-- Vector3 / number: Vector3
-- number / Vector3: Vector3
--@doc: Dividing two vectors returns a Vector3 with each component of the first
-- divided by those of the second.
--
-- A vector may also be divided by a number, which returns a Vector3 with each
-- component divided by the number.
function Vector3.__div(a, b)
	if type(a) == "number" then
		return newVector3(a/b.X, a/b.Y, a/b.Z)
	elseif type(b) == "number" then
		return newVector3(a.X/b, a.Y/b, a.Z/b)
	end
	return newVector3(a.X/b.X, a.Y/b.Y, a.Z/b.Z)
end

--@sec: Negation {Vector3.__unm}
--@ord: 1
--@def: -Vector3: Vector3
--@doc: Negating a vector returns a Vector3 with each the negation of each of
-- component.
function Vector3:__unm()
	return newVector3(-self.X, -self.Y, -self.Z)
end

return {
	--@sec: Vector3.new
	--@ord: -2
	--@def: Vector3.new(x: number?, y: number?, z: number?): Vector3
	--@doc: new returns a new Vector3 with the given components.
	new = newVector3,
}
