import { VehicleType } from "@/type/enum/vehicle-type";
import VehiclePage from "../page";

export default function NonTrontonPage() {
  return (
    <VehiclePage
      queryKey="vehicles-non-tronton"
      types={[VehicleType.TRUCK, VehicleType.PICKUP]}
    />
  );
}
