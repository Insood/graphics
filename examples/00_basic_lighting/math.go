package main

import "math"

type Point2 struct {
	x float64
	y float64
}

type Point2Int struct {
	x int
	y int
}

type Vector2Int struct {
	x int
	y int
}

type Vector2 struct {
	x float64
	y float64
}

func (v1 Vector2) cross(v2 Vector2) float64 {
	return v1.x*v2.y - v1.y*v2.x
}

type Point3 struct {
	x float64
	y float64
	z float64
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Color3 struct {
	r float64
	g float64
	b float64
}

func (v1 Vector3) cross(v2 Vector3) Vector3 {
	return Vector3{
		v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x,
	}
}

func (v Vector3) magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vector3) normalize() Vector3 {
	return Vector3{
		v.x / v.magnitude(),
		v.y / v.magnitude(),
		v.z / v.magnitude(),
	}
}

func (v1 Vector3) dot(v2 Vector3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v1 Vector2Int) dot(v2 Vector2Int) int {
	return v1.x*v2.x + v1.y*v2.y
}

func (v1 Point2) round() Point2Int {
	return Point2Int{
		int(math.Round(v1.x)),
		int(math.Round(v1.y)),
	}
}

func (v1 Point3) subtract(v2 Point3) Vector3 {
	return Vector3{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func (v1 Point2) subtract(v2 Point2) Vector2 {
	return Vector2{
		v1.x - v2.x,
		v1.y - v2.y,
	}
}

func (v1 Point2Int) subtract(v2 Point2Int) Vector2Int {
	return Vector2Int{
		v1.x - v2.x,
		v1.y - v2.y,
	}
}

// // return a new vector representing p1-p2
// float[] subtract(float[] p1, float p2[])
// {
//   float[] result = new float[p1.length];
//   for(int i = 0; i < p1.length; i++){
//      result[i] = p1[i] - p2[i];
//   }

//   return result;
// }

// int[] subtract(int[] p1, int p2[])
// {
//   int[] result = new int[p1.length];
//   for(int i = 0; i < p1.length; i++){
//      result[i] = p1[i] - p2[i];
//   }

//   return result;
// }

// float[] add(float[] p1, float[] p2) {
//     float[] result = new float[p1.length];
//   for(int i = 0; i < p1.length; i++){
//      result[i] = p1[i] + p2[i];
//   }

//   return result;
// }

// int[] round(float[] arr){
//   int[] result = new int[arr.length];
//   for(int i =0; i< arr.length; i++){
//     result[i] = int(arr[i]);
//   }
//   return result;
// }

// float[] multiply(float[] arr, float scalar){
//   float[] result = new float[arr.length];
//   for(int i =0; i < arr.length; i++){
//     result[i] = arr[i]*scalar;
//   }

//   return result;
// }
