// Lets TypeScript type `import schema from "./schema.sql" with { type: "text" }`
declare module "*.sql" {
	const content: string;
	export default content;
}
