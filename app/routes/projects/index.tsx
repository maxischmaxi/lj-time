import { getProjects } from "~/api";
import type { Route } from "./+types";

export async function loader() {
    const projects = await getProjects();
    return { projects };
}

export default function Projects({
    loaderData: { projects },
}: Route.ComponentProps) {
    return (
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                </tr>
            </thead>
            <tbody>
                {projects.map((project) => (
                    <tr key={project.id}>
                        <td>{project.id}</td>
                        <td>{project.name}</td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}
