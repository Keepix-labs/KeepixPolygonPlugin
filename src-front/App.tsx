import { Routes, Route } from "react-router-dom";
import HomePage from "./pages/Home";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

export default function App() {
  const queryClient = new QueryClient();

  return (
    <div className="App">
      <QueryClientProvider client={queryClient}>
        <Routes>
          <Route path={process.env.PUBLIC_URL + "/"} element={<HomePage />} />
        </Routes>
      </QueryClientProvider>
    </div>
  );
}
