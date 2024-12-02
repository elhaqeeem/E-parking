CREATE TABLE "parking_spots" (
  "id" int NOT NULL AUTO_INCREMENT,
  "spot_number" varchar(10) NOT NULL,
  "is_occupied" tinyint(1) DEFAULT '0',
  PRIMARY KEY ("id")
) 

CREATE TABLE "reservations" (
  "id" int NOT NULL AUTO_INCREMENT,
  "name" varchar(100) NOT NULL,
  "car_number" varchar(20) NOT NULL,
  "spot_id" int NOT NULL,
  "start_time" datetime NOT NULL,
  "duration" int NOT NULL,
  PRIMARY KEY ("id"),
  KEY "spot_id" ("spot_id"),
  CONSTRAINT "reservations_ibfk_1" FOREIGN KEY ("spot_id") REFERENCES "parking_spots" ("id")
) 