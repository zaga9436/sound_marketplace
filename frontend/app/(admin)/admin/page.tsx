import { Badge } from "@/shared/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/card";

export default function AdminPage() {
  return (
    <div className="space-y-8">
      <div className="space-y-3">
        <Badge>Admin</Badge>
        <div className="space-y-2">
          <h1>Admin workspace foundation.</h1>
          <p>Здесь дальше пойдут moderation users, cards, disputes и audit trail views.</p>
        </div>
      </div>

      <Card className="border-white/60 bg-white/90">
        <CardHeader>
          <CardTitle>Что уже предусмотрено</CardTitle>
        </CardHeader>
        <CardContent>
          <p>
            Защита через middleware, role-aware navigation и отдельная route group под admin flow
            уже заложены.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
