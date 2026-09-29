export interface TokenStorage {
	getToken(): Promise<string | undefined>;
	setToken(token: string): Promise<void>;
}
