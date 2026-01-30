package utils

func CalculateNDFL(salary float64) float64 {
	if salary <= 0 {
		return 0
	}

	annualSalary := salary * 12

	if annualSalary <= 2400000 {
		return 13.0
	}
	if annualSalary <= 5000000 {
		return 15.0
	}
	if annualSalary <= 20000000 {
		return 18.0
	}
	if annualSalary <= 50000000 {
		return 20.0
	}
	return 22.0
}
