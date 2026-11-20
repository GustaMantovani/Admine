import type { PortProtocol } from "../../domain/port-purpose";
import type { PortProbe } from "../../ports/port-probe";

export class NetPortProbe implements PortProbe {
	async isFree(port: number, protocol: PortProtocol): Promise<boolean> {
		try {
			if (protocol === "tcp") {
				const listener = Bun.listen({
					hostname: "0.0.0.0",
					port,
					socket: { data() {} },
				});
				listener.stop(true);
			} else {
				const socket = await Bun.udpSocket({ hostname: "0.0.0.0", port });
				socket.close();
			}
			return true;
		} catch {
			return false;
		}
	}
}
