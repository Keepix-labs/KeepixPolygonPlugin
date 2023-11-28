import { Dispatch, SetStateAction } from "react";

export const getErrorMsg = (e: unknown) => {
  if (typeof e === "string") {
    return e;
  }

  if (e instanceof Error) {
    return e.message;
  }

  return "Error";
};

export const safeFetch = async (
  url: string
) => {
  let res = new Response();
  try {
    res = await fetch(url, {
      headers: {
        "Cache-Control": "no-cache",
      },
    });
  } catch (error) {
    console.error(error);
    await new Promise((r) => setTimeout(r, 3000));

    try {
      res = await fetch(url, {
        headers: {
          "Cache-Control": "no-cache",
        },
      });
    } catch (error) {
      throw new Error("API not reachable");
    }
  }

  return res;
};
