import { VehicleType } from "@/type/enum/vehicle-type";
import VehiclePage from "../page";

export default function TrontonPage() {
  return (
    <VehiclePage queryKey="vehicles-tronton" types={[VehicleType.TRONTON]} />
  );
}
