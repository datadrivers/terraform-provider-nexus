##########
#
# This Script generates the MD Files for Nexus Data Sources and Resource Objects if not present
# and populates the nexus.erb layout file accordingly
#
##########

#!/bin/bash

cat > nexus.erb << EOI
<% wrap_layout :inner do %>
  <% content_for :sidebar do %>
    <div class="docs-sidebar hidden-print affix-top" role="complementary">
      <ul class="nav docs-sidenav">
        <li<%= sidebar_current("docs-home") %>>
        <a href="/docs/providers/index.html">All Providers</a>
        </li>

        <li<%= sidebar_current("docs-nexus-index") %>>
        <a href="/docs/providers/nexus/index.html">Nexus Provider</a>
        </li>
        <li<%= sidebar_current("docs-nexus-datasource") %>>
          <a href="#">Data Sources</a>
          <ul class="nav nav-visible">
EOI

writingDataSources=true

for res in $(ls nexus/data_source_* nexus/resource_* | grep -v test | sort )
do 

  res=$(echo $res | awk -F/ '{print $2}' | sed 's/\.go//g')

  if [[ "$res" =~ "resource" ]]
  then
    typeUrl="resources"
    typeLong="resource"
    if [[ "${writingDataSources}" == true ]]
    then
      cat >> nexus.erb << EOI
          </ul>
        </li>
        <li<%= sidebar_current("docs-nexus-resource") %>>
        <a href="#">Resources</a>
        <ul class="nav nav-visible">
EOI
      writingDataSources=false
    fi
  else
    typeUrl="data-sources"
    typeLong="data"
  fi

  echo $res

  if [[ ! -f docs/$typeUrl/$res.md ]]
  then
    cat > docs/$typeUrl/${res}.md << EOI
---
layout: "nexus"
page_title: "Nexus: $res"
sidebar_current: "docs-nexus-$typeLong-source"
description: |-
  Sample $typeLong source in the Terraform provider Nexus.
---

# $res

Sample data source in the Terraform provider scaffolding.

## Example Usage

\`\`\`hcl
data "scaffolding_data_source" "example" {
  sample_attribute = "foo"
}
\`\`\`

## Attributes Reference

* sample_attribute - Sample attribute.

EOI
  fi

  cat >> nexus.erb << EOI
            <li<%= sidebar_current("$res") %>>
              <a href="/docs/providers/nexus/$typeUrl/${res}.html">$res</a>
            </li>
EOI
done

cat >> nexus.erb << EOI
        </ul>
        </li>
      </ul>
    </div>
  <% end %>

  <%= yield %>
  <% end %>
EOI
