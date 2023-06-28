import { useEffect } from 'react';
import { useResource } from '../../../custom/api';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../../utils';
import Intl from '../../../i18n/i18n';

import StorePlugin from './StorePlugin';
import { UseQueryResult } from '@tanstack/react-query';
import { sortPlugins } from '../../../utils/sortArray';
import { MassaPluginModel, MassaStoreModel } from '../../../models';

function StoreSection({
  getPlugins,
}: {
  getPlugins: UseQueryResult<MassaPluginModel[], undefined>;
}) {
  const navigate = useNavigate();
  const {
    error,
    data: plugins,
    isLoading,
  } = useResource<MassaStoreModel[]>('plugin-store');

  const { refetch, data: myPlugins } = getPlugins;

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  });

  const isDownloaded = (plugin: MassaStoreModel) => {
    return (
      myPlugins?.some((myPlugin) => {
        const { name, author } = myPlugin;
        return name === plugin.name && author === plugin.author;
      }) || false
    );
  };

  return (
    <>
      {isLoading ? (
        <div className="mas-body mb-4 text-neutral">
          {Intl.t('store.loading')}
        </div>
      ) : plugins && plugins.length ? (
        <div className="flex gap-4 flex-wrap">
          {sortPlugins(plugins).map(
            (plugin, index: number) =>
              !isDownloaded(plugin) && (
                <StorePlugin key={index} plugin={plugin} refetch={refetch} />
              ),
          )}
        </div>
      ) : (
        <></>
      )}
    </>
  );
}

export default StoreSection;
