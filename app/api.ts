import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import {
    CustomerService,
    type CreateCustomer,
    type Customer,
} from "customer/v1/customer_pb";
import { ProjectService, type Project } from "project/v1/project_pb";

const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
    httpVersion: "1.1",
});

const customerClient = createClient(CustomerService, transport);
const projectClient = createClient(ProjectService, transport);

export async function getCustomers(): Promise<Customer[]> {
    return await customerClient.getCustomers({}).then((res) => res.customers);
}

export async function getProjects(): Promise<Project[]> {
    return await projectClient.getProjects({}).then((res) => res.projects);
}

export async function getProjectsByCustomer(id: string): Promise<Project[]> {
    return await projectClient
        .getProjectsByCustomer({ customerId: id })
        .then((res) => res.projects);
}

export async function deleteProject(id: string): Promise<string> {
    return await projectClient.deleteProject({ id }).then((res) => res.id);
}

export async function updateProject(
    data: Project
): Promise<Project | undefined> {
    return await projectClient
        .updateProject({ project: data })
        .then((res) => res.project);
}

export async function getProject(id: string): Promise<Project | undefined> {
    return await projectClient.getProject({ id }).then((res) => res.project);
}

export async function getCustomer(id: string): Promise<Customer | undefined> {
    return await customerClient.getCustomer({ id }).then((res) => res.customer);
}

export async function deleteCustomer(id: string): Promise<string> {
    return await customerClient.deleteCustomer({ id }).then((res) => res.id);
}

export async function updateCustomer(
    data: Customer
): Promise<Customer | undefined> {
    return await customerClient
        .updateCustomer({ customer: data })
        .then((res) => res.customer);
}

export async function createCustomer(
    data: CreateCustomer
): Promise<Customer | undefined> {
    return await customerClient
        .createCustomer({ customer: data })
        .then((res) => res.customer);
}
