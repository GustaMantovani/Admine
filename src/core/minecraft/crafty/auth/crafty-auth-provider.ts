export interface CraftyAuthProvider {
	getToken(): Promise<string>;
}
