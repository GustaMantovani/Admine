import { AdmineError } from "../../../shared/errors";

export class ServerNotFound extends AdmineError {}
export class ServerNameTaken extends AdmineError {}
export class InvalidServerName extends AdmineError {}
export class InvalidCreateOption extends AdmineError {}
export class InvalidServerPorts extends AdmineError {}
export class UnsupportedGameType extends AdmineError {}
export class CapabilityNotSupported extends AdmineError {}
export class NotApplicable extends AdmineError {}
export class InvalidFile extends AdmineError {}
export class PortRangeExhausted extends AdmineError {}
export class SharedPortBusy extends AdmineError {}
export class OrchestratorActionError extends AdmineError {}
