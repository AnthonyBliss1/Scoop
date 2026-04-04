import {
  Collection,
  KV,
  Method,
  Response,
  Scoop,
  Server,
  SyncData,
} from "../../../bindings/changeme";
import { getContext, setContext } from "svelte";

interface AppState {
  currentCollection: Collection;
  currentScoop: Scoop;
  currentServer: Server;
  allScoops: Scoop[];

  method: Method;
  url: string;
  response: Response;
  curlCommand: string;

  headers: KV[];
  queryParams: KV[];
  body: string;
}

export class AppStateClass implements AppState {
  currentCollection: Collection = $state(new Collection({ name: "temp" }));
  currentScoop: Scoop = $state(new Scoop({ name: "temp" }));
  currentServer: Server = $state(new Server({ key: "", url: "" }));
  allScoops: Scoop[] = $state([]);
  syncData: SyncData = $state(new SyncData());

  method: Method = $state(Method.Empty);
  url: string = $state("");
  response: Response = $state(new Response());
  curlCommand: string = $state("");

  headers: KV[] = $state([]);
  queryParams: KV[] = $state([]);
  body: string = $state("");

  reset() {
    this.currentCollection = new Collection({ name: "temp" });
    this.currentScoop = new Scoop({ name: "temp" });
    this.allScoops = [];
    this.syncData = new SyncData({ versionNum: 0 });

    this.method = Method.Empty;
    this.url = "";
    this.response = new Response();
    this.curlCommand = "";

    this.headers = [];
    this.queryParams = [];
    this.body = "";
  }

  resetResponse() {
    this.response = new Response();
  }
}

const DEFAULT_KEY = "$_app_state";

export const getAppState = (key = DEFAULT_KEY) => {
  return getContext<AppStateClass>(key);
};

export const setAppState = (key = DEFAULT_KEY) => {
  const appState = new AppStateClass();
  return setContext(key, appState);
};
