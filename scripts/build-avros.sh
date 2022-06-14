#!/bin/sh

gogen-avro ./api/avro/avro_gencode/ \
./api/avro/schemas/organization.avsc \
./api/avro/schemas/school.avsc \
./api/avro/schemas/user.avsc \
./api/avro/schemas/class.avsc \
./api/avro/schemas/organization_membership.avsc \
./api/avro/schemas/school_membership.avsc \
./api/avro/schemas/s3filecreated.avsc
