enum VehicleType {
  PICKUP = "PICKUP",
  TRONTON = "TRONTON",
  TRUCK = "TRUCK",
}

type Vehicle = {
  id: number;
  plate: string;
  type: VehicleType;
};
