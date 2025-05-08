import { getCustomers } from "~/api";
import type { Route } from "./+types";
import { CustomersTable } from "~/components/customers-table";
import type { Customer } from "customer/v1/customer_pb";

export async function loader() {
    const customers = await getCustomers();
    return { customers };
}

export default async function Customers({
    loaderData: { customers },
}: Route.ComponentProps) {
    return (
        <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
            {customers && (
                <CustomersTable customers={customers as Customer[]} />
            )}
        </div>
    );
}
