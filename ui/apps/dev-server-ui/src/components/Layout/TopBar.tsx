import { InngestLogo } from '@inngest/components/icons/logos/InngestLogo';
import { Link } from '@tanstack/react-router';

import { useInfoQuery } from '@/store/devApi';
import { SettingsMenu } from '../NavigationV2/SettingsMenu';

export default function TopBar() {
  const { data: info } = useInfoQuery();
  // Only a real dev server (`inngest dev`) exposes the /dev info endpoint with
  // isSingleNodeService falsy. Self-hosted `inngest start` runs in cloud mode,
  // where /dev 404s and `info` never loads — so we render "SERVER". Gating on a
  // positively-loaded `info` also avoids flashing "DEVELOPMENT SERVER" on load.
  const isDevServer = Boolean(info && !info.isSingleNodeService);

  return (
    <header className="bg-canvasSubtle relative z-[60] flex h-[48px] shrink-0 items-center justify-between gap-3 px-3">
      <div className="flex h-8 items-center gap-1.5">
        <Link to="/">
          <InngestLogo className="text-basis" width={96} />
        </Link>
        <span className="text-primary-intense text-[11px] font-medium leading-none">
          {isDevServer ? 'DEVELOPMENT SERVER' : 'SERVER'}
        </span>
      </div>
      <div className="flex items-center gap-3">
        <SettingsMenu />
      </div>
    </header>
  );
}
