CREATE TYPE "VehicleType" AS ENUM ('TRONTON', 'TRUCK', 'PICKUP');
CREATE TYPE "VehicleHistoryType" AS ENUM ('INCOME', 'EXPENSE');


CREATE TABLE "vehicles" (
    "id" SERIAL PRIMARY KEY,
    "plate" VARCHAR(100) NOT NULL,
    "type" "VehicleType" NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);

CREATE TABLE "vehicle_history" (
    "id" SERIAL PRIMARY KEY,
    "date" TIMESTAMP(3) NOT NULL,
    "description" TEXT NOT NULL,
    "type" "VehicleHistoryType" NOT NULL,
    "amount" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "profit" INTEGER ,
    "sack" INTEGER ,
    "vehicle_id" INTEGER NOT NULL REFERENCES "vehicles"("id") ON DELETE RESTRICT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(3)
);