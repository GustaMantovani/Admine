export interface RconPasswordStorage {
	getPassword(): Promise<string | undefined>;
	setPassword(password: string): Promise<void>;
}
