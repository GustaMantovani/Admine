export class UnknownMinecraftServerType extends Error {}

export class MinecraftServerActionError extends Error {}

export class CraftyApiError extends Error {
    constructor(
        public readonly status: number,
        public readonly body: string
    ) {
        super(
            `Crafty API responded with ${status}: ${CraftyApiError.summarize(body)}`
        );
    }

    // Crafty answers with { error, error_data } JSON or, on unhandled errors, a Python traceback
    private static summarize(body: string): string {
        try {
            const json = JSON.parse(body) as { error?: string; error_data?: string; };
            if (json.error !== undefined) {
                return json.error_data
                    ? `${json.error} (${json.error_data})`
                    : json.error;
            }
        } catch { }

        const lines = body.trim().split("\n");
        return lines[lines.length - 1] ?? body;
    }
}
