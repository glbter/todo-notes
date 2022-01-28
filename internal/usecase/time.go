package usecase

import (
	"time"
	"todoNote/internal/model"
)

var (
	hour = 60*60
	UTCp1  = time.FixedZone(model.UTCp1, 1*hour)
	UTCp2  = time.FixedZone(model.UTCp2, 2*hour)
	UTCp3  = time.FixedZone(model.UTCp3, 3*hour)
	UTCp4  = time.FixedZone(model.UTCp4, 4*hour)
	UTCp5  = time.FixedZone(model.UTCp5, 5*hour)
	UTCp6  = time.FixedZone(model.UTCp6, 6*hour)
	UTCp7  = time.FixedZone(model.UTCp7, 7*hour)
	UTCp8  = time.FixedZone(model.UTCp8, 8*hour)
	UTCp9  = time.FixedZone(model.UTCp9, 9*hour)
	UTCp10  = time.FixedZone(model.UTCp10, 10*hour)
	UTCp11  = time.FixedZone(model.UTCp11, 11*hour)
	UTCp12  = time.FixedZone(model.UTCp12, 12*hour)

	UTCm1  = time.FixedZone(model.UTCm1, -1*hour)
	UTCm2  = time.FixedZone(model.UTCm2, -2*hour)
	UTCm3  = time.FixedZone(model.UTCm3, -3*hour)
	UTCm4  = time.FixedZone(model.UTCm4, -4*hour)
	UTCm5  = time.FixedZone(model.UTCm5, -5*hour)
	UTCm6  = time.FixedZone(model.UTCm6, -6*hour)
	UTCm7  = time.FixedZone(model.UTCm7, -7*hour)
	UTCm8  = time.FixedZone(model.UTCm8, -8*hour)
	UTCm9  = time.FixedZone(model.UTCm9, -9*hour)
	UTCm10  = time.FixedZone(model.UTCm10, -10*hour)
	UTCm11  = time.FixedZone(model.UTCm11, -11*hour)
	UTCm12  = time.FixedZone(model.UTCm12, -12*hour)

	converter = make(map[model.TimeZone] *time.Location)
)


func init()  {
	converter[model.UTC] = time.UTC
	converter[model.UTCp1] = UTCp1
	converter[model.UTCp2] = UTCp2
	converter[model.UTCp3] = UTCp3
	converter[model.UTCp4] = UTCp4
	converter[model.UTCp5] = UTCp5
	converter[model.UTCp6] = UTCp6
	converter[model.UTCp7] = UTCp7
	converter[model.UTCp8] = UTCp8
	converter[model.UTCp9] = UTCp9
	converter[model.UTCp10] = UTCp10
	converter[model.UTCp11] = UTCp11
	converter[model.UTCp12] = UTCp12

	converter[model.UTCm1] = UTCm1
	converter[model.UTCm2] = UTCm2
	converter[model.UTCm3] = UTCm3
	converter[model.UTCm4] = UTCm4
	converter[model.UTCm5] = UTCm5
	converter[model.UTCm6] = UTCm6
	converter[model.UTCm7] = UTCm7
	converter[model.UTCm8] = UTCm8
	converter[model.UTCm9] = UTCm9
	converter[model.UTCm10] = UTCm10
	converter[model.UTCm11] = UTCm11
	converter[model.UTCm12] = UTCm12
}

func Convert(dateTime time.Time, zone model.TimeZone) time.Time {
	loc := converter[zone]
	return dateTime.In(loc)
}

func ValidateZone(z string) (model.TimeZone, bool) {
	_, ok := converter[z]
	return z, ok
}
