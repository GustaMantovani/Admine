export interface CraftyAuthProvider {
	getToken(): Promise<string>;
	getNewToken(): Promise<string>;
}
