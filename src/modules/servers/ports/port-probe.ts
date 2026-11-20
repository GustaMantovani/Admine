import type { PortProtocol } from "../domain/port-purpose";

export interface PortProbe {
	isFree(port: number, protocol: PortProtocol): Promise<boolean>;
}
