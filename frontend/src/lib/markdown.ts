import { marked } from "marked";
import DOMPurify from "isomorphic-dompurify";
export function renderMarkdown(source: string): string {
  return DOMPurify.sanitize(marked.parse(source, { async: false }) as string);
}
