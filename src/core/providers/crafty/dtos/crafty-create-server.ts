import type { CraftyCreateBedrockServer } from "./crafty-create-bedrock-server";
import type { CraftyCreateHytaleServer } from "./crafty-create-hytale-server";
import type { CraftyCreateJavaServer } from "./crafty-create-java-server";

export type CraftyCreateServer = CraftyCreateJavaServer | CraftyCreateBedrockServer | CraftyCreateHytaleServer;
