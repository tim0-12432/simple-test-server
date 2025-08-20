import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function InputWithHelperTextDemo() {
  return (
    <div className="w-full max-w-xs space-y-1.5">
      <Label htmlFor="email-address">Email Address</Label>
      <Input id="email-address" type="email" placeholder="Email" />
      <p className="text-[0.8rem] text-muted-foreground">
        We&apos;ll never share your email with anyone else.
      </p>
    </div>
  );
}
