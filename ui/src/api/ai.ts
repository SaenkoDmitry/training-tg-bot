import {api} from "./client";

export const getAIProgramContext = (programId?: number) => {
    const query = programId ? `?program_id=${programId}` : "";
    return api<AIProgramContext>(`/api/ai/program-context${query}`);
};

export const buildAIProgramPrompt = (request: AIProgramPromptRequest) =>
    api<AIProgramPromptResponse>("/api/ai/program-prompt", {
        method: "POST",
        body: JSON.stringify(request),
    });
