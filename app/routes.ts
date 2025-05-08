import {
    type RouteConfig,
    index,
    prefix,
    route,
} from "@react-router/dev/routes";

export default [
    index("routes/index.tsx"),
    ...prefix("/projects", [
        index("routes/projects/index.tsx"),
        route("/create", "routes/projects/create.tsx"),
        route("/edit/:projectId", "routes/projects/edit.tsx"),
    ]),
    ...prefix("/customers", [
        index("routes/customers/index.tsx"),
        route("/create", "routes/customers/create.tsx"),
        route("/edit/:customerId", "routes/customers/edit.tsx"),
    ]),
] satisfies RouteConfig;
