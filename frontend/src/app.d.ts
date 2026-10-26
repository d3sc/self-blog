declare global {
  namespace App {
    interface Locals {
      admin?: Admin;
    }
  }
}
import type { Admin } from "$lib/types";
export {};
