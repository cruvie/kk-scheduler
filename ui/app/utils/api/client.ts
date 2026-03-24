import {createClient} from "@connectrpc/connect";

import {transport} from "./api";
import {KKSchedule} from "~~/gen/kk_schedule/service_pb";

export const clientKKSchedule = createClient(KKSchedule, transport);
