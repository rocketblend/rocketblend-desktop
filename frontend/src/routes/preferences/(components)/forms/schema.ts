import { z } from "zod";

export const desktopGeneralForm = z.object({
    watchFolder: z.string().regex(/^([a-zA-Z]:)?(\\[a-zA-Z0-9_-]+)+\\?$/, { message: "Invalid file path format." }),
});

export type DesktopGeneralForm = typeof desktopGeneralForm;