import {
    type RouteConfig,
    index,
    layout,
    prefix,
    route,
} from "@react-router/dev/routes";

export default [
    layout("layout.tsx", [
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
    ]),
] satisfies RouteConfig;
