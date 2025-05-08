import { getCustomers } from "~/api";
import type { Route } from "./+types";
import { Link } from "react-router";

export async function loader() {
    const customers = await getCustomers();
    return { customers };
}

export default async function Customers({
    loaderData: { customers },
}: Route.ComponentProps) {
    return (
        <div>
            <Link to="/customers/create">Create Customer</Link>
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Email</th>
                        <th>Phone</th>
                    </tr>
                </thead>
                <tbody>
                    {customers?.map((customer) => (
                        <tr key={customer.id}>
                            <td>{customer.id}</td>
                            <td>{customer.name}</td>
                            <td>{customer.email}</td>
                            <td>{customer.phone}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
