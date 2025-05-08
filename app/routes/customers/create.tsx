import { getCustomers } from "~/api";
import type { Route } from "../+types";

export async function loader() {
    const customers = await getCustomers();
    return { customers };
}
export default function CreateCustomer({
    loaderData: { customers },
}: Route.ComponentProps) {
    return (
        <div>
            <h1>Create Customer</h1>
            <form></form>
        </div>
    );
}
